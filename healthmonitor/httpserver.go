package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	router := gin.Default()
	router.LoadHTMLGlob("resources/*.tmp")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmp", nil)
	})

	router.GET("/results", func(c *gin.Context) {
		start, _ := strconv.Atoi(c.Query("start"))
		end, _ := strconv.Atoi(c.Query("end"))

		if start < 1 || end > 5 || start > end {
			c.HTML(http.StatusOK, "error.tmp", gin.H{
				"error": "Invalid Parameters",
			})
			return
		}

		c.HTML(http.StatusOK, "results.tmp", gin.H{
			"result": template.HTML(runMapRedJob(start, end)),
		})
	})

	router.Run(":5000")
}

func getFilesList(start int, end int) string {
	str := ""

	for i := start; i <= end; i++ {
		str += fmt.Sprintf("/myfile_%v.json", i)

		if i != end {
			str += ","
		}
	}

	return str
}

func runMapRedJob(start int, end int) string {
	clearHDFSOutput()

	hadoopPath := "/home/hadoop/hadoop/bin/hadoop"
	jarPath := "/home/hadoop/MapReduceHealthMessages.jar"
	filesList := getFilesList(start, end)
	fmt.Println("Files list:", filesList)

	cmd := exec.Command(hadoopPath, "jar", jarPath, filesList, "/output")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Running Job: ")
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Command finished with error: %v", err)
	}

	fmt.Println("Job finished")

	return formatOutput(readHDFSOutput())
	// return readHDFSOutput()
}

func readHDFSOutput() string {
	out, err := exec.Command("/home/hadoop/hadoop/bin/hdfs", "dfs", "-cat", "/output/part-00000").Output()
	if err != nil {
		fmt.Printf("Failed to read /output: %v", err)
	}

	return string(out)
}

func formatOutput(hdfsOutput string) string {
	fmt.Println("hdfs: ", hdfsOutput)
	services := strings.Fields(hdfsOutput)
	output := ""
	fmt.Println("services", services)
	fmt.Println(len(services))

	for i := 0; i < len(services); i += 2 {
		fmt.Println(i)
		service := services[i : i+2]
		fmt.Println(service)
		serviceName := service[0]
		fmt.Println(serviceName)
		metrics := strings.Split(service[1], ",")
		fmt.Println("Metrics: ", len(metrics))

		output += fmt.Sprintf("<h1>Service Name: %v</h1>", serviceName)
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Avg CPU: %v<br>", metrics[0])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Avg RAM Total: %v<br>", metrics[1])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Avg RAM Free: %v<br>", metrics[2])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Avg Disk Total: %v<br>", metrics[3])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Avg Disk Free: %v<br>", metrics[4])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Peak CPU: %v<br>", metrics[5])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Peak RAM Total: %v<br>", metrics[6])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Peak RAM Free: %v<br>", metrics[7])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Peak Disk Total: %v<br>", metrics[8])
		output += fmt.Sprintf("&nbsp;&nbsp;&nbsp;&nbsp;Peak Disk Free: %v<br>", metrics[9])

		output += "<hr>"
	}

	return output

}

func clearHDFSOutput() {
	fmt.Print("Deleting /output:")

	cmd := exec.Command("/home/hadoop/hadoop/bin/hdfs", "dfs", "-rm", "-r", "/output")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	} else {
		fmt.Println("Success")
	}

}
