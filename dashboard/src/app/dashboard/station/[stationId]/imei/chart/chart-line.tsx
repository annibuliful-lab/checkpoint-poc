// ----------------------------------------------------------------------

import Chart, { useChart } from "@/components/chart";

type Props = {
  series: {
    name: string;
    data: number[];
  }[];
};

export default function ImeiChartLine({ series }: Props) {
  const chartOptions = useChart({
    xaxis: {
      categories: [
        "Jan",
        "Feb",
        "Mar",
        "Apr",
        "May",
        "Jun",
        "Jul",
        "Aug",
        "Sep",
      ],
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
