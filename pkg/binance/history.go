package binance

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"sort"
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
	for _, acc := range Accounts {
		balance, err := GetFuturesBalance(acc.AccountType, cfg, acc.Email)
		if err != nil {
			return "", err
		}

		updatedHistories, err := upsertTodayBalance(acc.Label, balance, bucketName, fileName)
		if err != nil {
			return "", err
		}

		histories := filterHistoriesForAccount(updatedHistories, acc.Label)
		for _, history := range histories {
			log.Printf("Date: %s, Account: %s, Balance: %f", history.Date, history.Account, history.Balance)
		}

		calculatedReturns, err := calculateReturns(histories, acc.Label)
		if err != nil {
			return "", err
		}

		totalResults.WriteString(calculatedReturns)
	}

	return totalResults.String(), nil
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

		builder.WriteString("Date,Account,Balance\n")
		for _, acc := range Accounts {
			builder.WriteString(fmt.Sprintf("%s,%s,0\n", time.Now().Format("2006-01-02"), acc.Label))
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
	today := time.Now().Truncate(24 * time.Hour)

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

func calculateReturns(histories []BalanceHistory, accLabel string) (string, error) {
	if len(histories) == 0 {
		return "", fmt.Errorf("no history found for account: %s", accLabel)
	}

	sort.Slice(histories, func(i, j int) bool {
		return histories[i].Date.Before(histories[j].Date)
	})

	initialBalance := histories[0].Balance
	currentBalance := histories[len(histories)-1].Balance

	var overallReturn float64
	if initialBalance != 0 {
		overallReturn = (currentBalance - initialBalance) / initialBalance * 100.0
	} else {
		overallReturn = 0
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	dailyReturn := calculateReturnFromDaysAgo(histories, currentBalance, yesterday, accLabel)

	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	weeklyReturn := calculateReturnFromDaysAgo(histories, currentBalance, oneWeekAgo, accLabel)

	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	monthlyReturn := calculateReturnFromDaysAgo(histories, currentBalance, oneMonthAgo, accLabel)

	dailyReturnStr := formatReturn(dailyReturn)
	weeklyReturnStr := formatReturn(weeklyReturn)
	monthlyReturnStr := formatReturn(monthlyReturn)

	var resultBuilder strings.Builder
	resultBuilder.WriteString(":bank: " + accLabel + "\n")
	resultBuilder.WriteString(strings.Repeat("-", 40) + "\n")
	resultBuilder.WriteString(fmt.Sprintf("Total: %s (%.0f → %.0f)\n", formatReturn(overallReturn), initialBalance, currentBalance))
	resultBuilder.WriteString(fmt.Sprintf("1D: %s\n", dailyReturnStr))
	resultBuilder.WriteString(fmt.Sprintf("1W: %s\n", weeklyReturnStr))
	resultBuilder.WriteString(fmt.Sprintf("1M: %s\n", monthlyReturnStr))
	resultBuilder.WriteString(strings.Repeat("-", 40) + "\n")
	resultBuilder.WriteString("\n\n")

	return resultBuilder.String(), nil
}

func calculateReturnFromDaysAgo(histories []BalanceHistory, currentBalance float64, daysAgo time.Time, accountLabel string) float64 {
	for _, history := range histories {
		if history.Date.Year() == daysAgo.Year() && history.Date.Month() == daysAgo.Month() && history.Date.Day() == daysAgo.Day() {
			if history.Balance == 0 {
				return 0
			}
			return (currentBalance - history.Balance) * 100.0 / history.Balance
		}
	}
	return 0
}

func formatReturn(val float64) string {
	if val == 0 {
		return "-"
	}
	return fmt.Sprintf("%.0f%%", val)
}
