package main

import (
	"fmt"
	"gin-server/app/module/fine_tuning/model"
	"gin-server/app/module/fine_tuning/service"
	"math"
	"path/filepath"
	"time"
)

func main() {
	absPath, err := filepath.Abs(model.FilePath)
	if err != nil {
		panic(err)
	}

	var file *model.OpenAIFile
	file, _ = service.UploadFile(absPath, "fine-tune")
	fmt.Println(file)

	var job *model.FineTuningJobResponse
	job, _ = service.CreatedJob(file.ID)
	fmt.Println(job)

	// 指数退避
	var jobDetail *model.FineTuningJobResponse
	for job.Status != "succeeded" && job.Status != "failed" {
		fmt.Println("Waiting for job to complete...")

		for i := 0; i < 3; i++ {
			jobDetail, _ = service.GetJob(job.ID)
			if jobDetail.Status == "succeeded" {
				fmt.Printf("Job succeeded: (model: %s)", jobDetail.FineTunedModel)
				break
			}
			if jobDetail.Status == "failed" {
				fmt.Printf("Job failed: %s", jobDetail.Status)
				return
			}

			delay := time.Duration(math.Pow(2, float64(i))) * model.BaseDelay
			time.Sleep(delay)
		}
	}
}
