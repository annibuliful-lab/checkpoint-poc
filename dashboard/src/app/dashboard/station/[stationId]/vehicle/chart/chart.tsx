"use client";
import { Box, Button, Card, Stack, Typography } from "@mui/material";
import Chart from "chart.js/auto";
import VehicleHeader from "../filter";
import VehicleChartLine from "./chart-line";
Chart.register();
type Props = {
  actions?: JSX.Element[];
};
export default function VehicleChart({ actions }: Props) {
  return (
    <Stack spacing={1}>
      <VehicleHeader />
      <Card sx={{ p: 2 }}>
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

        <VehicleChartLine
          series={[
            {
              name: "Total",
              data: [10, 41, 35, 51, 49, 62, 69, 91, 148],
            },
          ]}
        />
      </Card>
    </Stack>
  );
}
