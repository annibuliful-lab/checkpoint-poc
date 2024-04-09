"use client";
import { Box, Button, Card, Stack, Typography } from "@mui/material";
import { Line } from "react-chartjs-2";
import Chart from "chart.js/auto";
import VehicleHeader from "./filter";
Chart.register();
type Props = {
  actions?: JSX.Element[];
};
export default function VehicleChart({ actions }: Props) {
  const labels = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];
  const data = {
    labels: labels,
    datasets: [
      {
        label: "My First Dataset",
        data: [65, 59, 80, 81, 56, 55, 40, 80, 81, 56, 55, 40],
        fill: false,
        borderColor: "rgb(75, 192, 192)",
        tension: 0.1,
      },
    ],
  };

  const options = {
    plugins: {
      legend: {
        display: false,
      },
    },
    scales: {
      x: {
        offset: true,
      },
      y: {},
    },
  };
  return (
    <Stack spacing={1}>
      <VehicleHeader />
      <Card sx={{ p: 2, height: 240 }}>
        <Stack direction={"row"} spacing={1} justifyContent={"space-between"}>
          <Stack direction={"row"} spacing={1} alignItems={"center"}>
            <Box sx={{ height: 18, width: 4, background: "#000000" }} />
            <Typography variant="subtitle2" fontWeight={700}>
              Overview
            </Typography>
            <Box sx={{ width: 10 }} />
            {["ALL", "IMSI", "IMEI"].map((filterType) => (
              <Button key={filterType} size="small" variant="outlined">
                {filterType}
              </Button>
            ))}
          </Stack>

          <Stack direction={"row"} spacing={1} alignItems={"center"}>
            {["DAY", "WEEK", "MONTH", "YEAR", "CUSTOM"].map((filterDate) => (
              <Button key={filterDate} size="small" variant="outlined">
                {filterDate}
              </Button>
            ))}
          </Stack>
        </Stack>
        {/* <Line data={data} options={options} /> */}
      </Card>
    </Stack>
  );
}
