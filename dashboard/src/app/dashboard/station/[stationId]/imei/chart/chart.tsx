"use client";
import { Box, Button, Card, Stack, Typography } from "@mui/material";
import Chart from "chart.js/auto";
import ImeiHeader from "../filter";
import ImeiChartLine from "./chart-line";
import {
  GetStationImeiImsiActivitySummaryCategory,
  useGetStationImeiImsiActivitySummaryQuery,
} from "@/apollo-client";
Chart.register();
type Props = {
  stationLocationId: string;
};
export default function ImeiChart({ stationLocationId }: Props) {
  const { loading, data } = useGetStationImeiImsiActivitySummaryQuery({
    variables: {
      stationId: stationLocationId,
      groupBy: GetStationImeiImsiActivitySummaryCategory.Day,
    },
  });
  return (
    <Stack spacing={1}>
      <ImeiHeader />
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

        <ImeiChartLine
          categories={data?.getStationImeiImsiActivitySummary.categories ?? []}
          series={data?.getStationImeiImsiActivitySummary.series ?? []}
        />
      </Card>
    </Stack>
  );
}
