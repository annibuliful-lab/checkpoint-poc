// ----------------------------------------------------------------------

import { StationImeiImsiActivitySummarySerie } from "@/apollo-client";
import Chart, { useChart } from "@/components/chart";

type Props = {
  categories: string[];
  series: StationImeiImsiActivitySummarySerie[];
};

export default function ImeiChartLine({ series, categories }: Props) {
  const chartOptions = useChart({
    xaxis: {
      categories,
    },
    tooltip: {
      x: {
        show: false,
      },
      marker: { show: false },
    },
  });

  return (
    <Chart
      dir="ltr"
      type="line"
      series={series.map((s) => ({ name: s.label, data: s.data }))}
      options={chartOptions}
      width="100%"
      height={240}
    />
  );
}
