// ----------------------------------------------------------------------

import { StationVehicleActivitySummarySerie } from "@/apollo-client";
import Chart, { useChart } from "@/components/chart";

type Props = {
  categories: string[];
  series: StationVehicleActivitySummarySerie[];
};

export default function VehicleChartLine({ series, categories }: Props) {
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
      series={series}
      options={chartOptions}
      width="100%"
      height={240}
    />
  );
}
