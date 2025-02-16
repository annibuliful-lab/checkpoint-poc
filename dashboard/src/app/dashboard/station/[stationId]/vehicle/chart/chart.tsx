"use client";
import { Box, Button, Card, Stack, Typography } from "@mui/material";
import Chart from "chart.js/auto";
import VehicleHeader from "../filter";
import VehicleChartLine from "./chart-line";
import {
  GetStationVehicleActivitySummaryCategory,
  useGetStationVehicleActivitySummaryQuery,
} from "@/apollo-client";
Chart.register();
type Props = {
  stationLocationId: string;
};
export default function VehicleChart({ stationLocationId }: Props) {
  const { loading, data } = useGetStationVehicleActivitySummaryQuery({
    variables: {
      stationId: stationLocationId,
      groupBy: GetStationVehicleActivitySummaryCategory.Day,
    },
  });
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
          categories={data?.getStationVehicleActivitySummary?.categories ?? []}
          series={data?.getStationVehicleActivitySummary?.series ?? []}
        />
      </Card>
    </Stack>
  );
}
