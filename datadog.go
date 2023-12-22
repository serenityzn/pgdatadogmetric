// Create a tag configuration returns "Created" response

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

func datadogCreateMetricTag() error {
	body := datadogV2.MetricTagConfigurationCreateRequest{
		Data: datadogV2.MetricTagConfigurationCreateData{
			Type: datadogV2.METRICTAGCONFIGURATIONTYPE_MANAGE_TAGS,
			Id:   "MyTestPgMetric",
			Attributes: &datadogV2.MetricTagConfigurationCreateAttributes{
				Tags: []string{
					"app",
					"datacenter",
				},
				MetricType: datadogV2.METRICTAGCONFIGURATIONMETRICTYPES_COUNT,
			},
		},
	}
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)
	resp, r, err := api.CreateTagConfiguration(ctx, "ExampleMetricMyCount", body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.CreateTagConfiguration`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return err
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.CreateTagConfiguration`:\n%s\n", responseContent)

	return nil
}

func datadogSubmitMetric(count float64) error {
	body := datadogV2.MetricPayload{
		Series: []datadogV2.MetricSeries{
			{
				Metric: "postgresql.tpart_charge_select_count",
				Type:   datadogV2.METRICINTAKETYPE_COUNT.Ptr(),
				Points: []datadogV2.MetricPoint{
					{
						Timestamp: datadog.PtrInt64(time.Now().Unix()),
						Value:     datadog.PtrFloat64(count),
					},
				},

				Resources: []datadogV2.MetricResource{
					{
						Name: datadog.PtrString("dev-pg-master.hvprcl2dpn2ezcozamnevzvxxc.ax.internal.cloudapp.net"),
						Type: datadog.PtrString("host"),
					},
				},
			},
		},
	}
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)
	resp, r, err := api.SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return err
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", responseContent)

	return nil
}
