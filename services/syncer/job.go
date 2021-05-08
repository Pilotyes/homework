package syncer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"shop-api/internal/config"
	"shop-api/internal/storage"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

//JobService ...
type JobService struct {
	name    string
	config  *config.SyncConfig
	logger  *logrus.Entry
	storage storage.Storage
}

//NewJobService ...
func NewJobService(name string, storage storage.Storage, logger *logrus.Logger, syncConfig *config.SyncConfig) (*JobService, error) {
	if name == "" || storage == nil || logger == nil || syncConfig == nil {
		return nil, errors.New("Couldn't create job")
	}

	l := logger.WithFields(logrus.Fields{
		"job name": name,
	})

	return &JobService{
		name:    name,
		config:  syncConfig,
		logger:  l,
		storage: storage,
	}, nil
}

//Start ...
func (j *JobService) Start() {
	j.logger.Infof("Started sync job \"%s\"", j.name)
	for {
		now := time.Now().UTC()
		callTime := time.Date(now.Year(), now.Month(), now.Day(), j.config.Hours, j.config.Minutes, 0, 0, time.UTC)
		if callTime.Before(now) {
			callTime = callTime.Add(time.Hour * 24)
		}

		j.logger.Infoln("Next sync at", callTime.Local())
		duration := callTime.Sub(time.Now().UTC())

		time.Sleep(duration)
		go func() {
			if err := j.Sync(); err != nil {
				j.logger.Errorln(err)
			}
		}()
	}
}

//Sync ...
func (j *JobService) Sync() error {
	j.logger.Infoln("Syncing...")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	j.logger.Debugln("Created new HTTP client", client)

	req, err := http.NewRequest(http.MethodGet, j.config.UrlToFile, nil)
	if err != nil {
		return err
	}
	j.logger.Debugln("Created new HTTP request", req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Couln't fetch file %s: %s", j.config.UrlToFile, resp.Status)
	}

	r := csv.NewReader(resp.Body)
	defer resp.Body.Close()
	data, err := r.ReadAll()
	if err != nil {
		return err
	}

	items := j.storage.Items().GetItems()
	resultDiscounts := make([]int, len(items))

	for _, v := range data {
		if len(v) < 3 {
			j.logger.Errorln("Can't parse discount data", v)
			continue
		}

		if v[0] == "k" {
			continue
		}

		discount, err := strconv.Atoi(v[2])
		if err != nil {
			j.logger.Errorln("Can't parse discount value", v[2], err)
			continue
		}

		switch v[0] {
		case "category":
			for i, item := range items {
				if item.Category == v[1] {
					resultDiscounts[i] += discount
				}
			}
		case "item":
			articul, err := strconv.Atoi(v[1])
			if err != nil {
				j.logger.Errorln("Can't parse articul value", v[1], err)
				continue
			}

			for i, item := range items {
				if item.Articul == articul {
					resultDiscounts[i] += discount
					item.ProductOfDay = true
				}
			}
		case "-":
			for i := range items {
				resultDiscounts[i] += discount
			}
		default:
			continue
		}
	}

	for i := range items {
		discountPrice := *(items[i].OriginalPrice) * float64(100-resultDiscounts[i]) / float64(100)
		items[i].DiscountPrice = &discountPrice
	}

	j.logger.Infoln("Sync done")
	return nil
}
