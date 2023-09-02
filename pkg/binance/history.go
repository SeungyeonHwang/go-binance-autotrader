package binance

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/SeungyeonHwang/go-binance-autotrader/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BalanceHistory struct {
	Date    time.Time
	Account string
	Balance float64
}

func FetchAllHistory(cfg *config.Config, bucketName, fileName string) (string, error) {
	_, err := initHistoryInS3(bucketName, fileName)
	if err != nil {
		return "", fmt.Errorf("error fetching history from S3: %w", err)
	}

	var totalResults strings.Builder
	totalInitialBalance := 0.0
	totalCurrentBalance := 0.0

	for _, acc := range Accounts {
		balance, err := GetFuturesBalance(acc.AccountType, cfg, acc.Email)
		if err != nil {
			return "", err
		}
		totalCurrentBalance += float64(balance)

		updatedHistories, err := upsertTodayBalance(acc.Label, balance, bucketName, fileName)
		if err != nil {
			return "", err
		}

		histories := filterHistoriesForAccount(updatedHistories, acc.Label)
		monthHistories := filterHistoriesForThisMonth(histories)

		initialBalance := 0.0
		if len(monthHistories) > 0 {
			initialBalance = monthHistories[0].Balance
		}
		totalInitialBalance += initialBalance

		totalDelta := (float64(balance) - initialBalance) / initialBalance * 100

		totalResults.WriteString(fmt.Sprintf("%s [%d (%+.2f%%)]\n", acc.Label, int(balance), totalDelta))
		totalResults.WriteString(strings.Repeat("-", 40) + "\n")

		for i, history := range monthHistories {
			delta := "(+0%)"
			if i > 0 {
				prevBalance := monthHistories[i-1].Balance
				deltaValue := (history.Balance - prevBalance) / prevBalance * 100
				delta = fmt.Sprintf("(%+d%%)", int(math.Round(deltaValue)))
			}
			totalResults.WriteString(fmt.Sprintf("%02d.%02d: %d %s\n", history.Date.Month(), history.Date.Day(), int(history.Balance), delta))
		}

		totalResults.WriteString(strings.Repeat("-", 40) + "\n")
		totalResults.WriteString("\n")
	}

	overallTotalDelta := (totalCurrentBalance - totalInitialBalance) / totalInitialBalance * 100
	totalResults.WriteString(strings.Repeat("=", 40) + "\n")
	totalResults.WriteString(fmt.Sprintf("Total Balance: %d → %d\n", int(totalInitialBalance), int(totalCurrentBalance)))
	totalResults.WriteString(fmt.Sprintf("Total Profit: %+7.2f%%\n", overallTotalDelta))
	totalResults.WriteString(strings.Repeat("=", 40) + "\n")

	return totalResults.String(), nil
}

func filterHistoriesForThisMonth(histories []BalanceHistory) []BalanceHistory {
	var monthHistories []BalanceHistory
	currentMonth := time.Now().Month()

	for _, history := range histories {
		if history.Date.Month() == currentMonth {
			monthHistories = append(monthHistories, history)
		}
	}
	return monthHistories
}

func DBClear(bucketName, fileName string) error {
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return err
	}

	return nil
}

func initHistoryInS3(bucketName, fileName string) ([]BalanceHistory, error) {
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	_, err := s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		var builder strings.Builder
		var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

		builder.WriteString("Date,Account,Balance\n")
		for _, acc := range Accounts {
			builder.WriteString(fmt.Sprintf("%s,%s,0\n", time.Now().In(jst).Format("2006-01-02"), acc.Label))
		}
		fileContent := builder.String()

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileName),
			Body:   bytes.NewReader([]byte(fileContent)),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create new file in S3: %w", err)
		}
	}
	return nil, err
}

func upsertTodayBalance(accountLabel string, balance int, bucketName, fileName string) ([]BalanceHistory, error) {
	today := time.Now().Add(9 * time.Hour).Truncate(24 * time.Hour)

	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	getResp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	records, err := csv.NewReader(getResp.Body).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv from S3: %w", err)
	}

	found := false
	for i, record := range records {
		if record[0] == today.Format("2006-01-02") && record[1] == accountLabel {
			records[i][2] = fmt.Sprintf("%d", balance)
			found = true
			break
		}
	}

	if !found {
		records = append(records, []string{today.Format("2006-01-02"), accountLabel, fmt.Sprintf("%d", balance)})
	}

	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	err = csvWriter.WriteAll(records)
	if err != nil {
		return nil, fmt.Errorf("failed to write updated records to csv: %w", err)
	}

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to put object to S3: %w", err)
	}

	var histories []BalanceHistory
	for _, record := range records {
		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			continue
		}
		account := record[1]
		bal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			continue
		}
		histories = append(histories, BalanceHistory{Date: date, Account: account, Balance: bal})

	}

	return histories, nil
}

func filterHistoriesForAccount(allHistories []BalanceHistory, accountLabel string) []BalanceHistory {
	var filteredHistories []BalanceHistory
	for _, history := range allHistories {
		if history.Account == accountLabel {
			filteredHistories = append(filteredHistories, history)
		}
	}
	return filteredHistories
}
