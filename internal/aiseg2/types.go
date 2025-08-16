package aiseg2

// PowerSummary は電力サマリー情報を保持します
type PowerSummary struct {
	TotalGenerationPowerKW MetricsElement   `json:"totalGenerationPowerKW"`
	TotalUsagePowerKW      MetricsElement   `json:"totalUsagePowerKW"`
	TotalBalancePowerKW    MetricsElement   `json:"totalBalancePowerKW"`
	DetailsGenerationPower []MetricsElement `json:"detailsGenerationPower"`
}

// DetailUsagePower は詳細消費電力情報を保持します
type DetailUsagePower []MetricsElement

// MetricsElement は個別のメトリクス要素を保持します
type MetricsElement struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
