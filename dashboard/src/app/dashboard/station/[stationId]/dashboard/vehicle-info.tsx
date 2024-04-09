import { Stack, Typography, Button, Divider, Box } from "@mui/material";
import { useAtomValue } from "jotai";
import React from "react";
import { stationDashboardActivityAtom } from "./store";
import { fDateTime } from "@/utils/format-time";

export default function VehicleInfo() {
  const dashboardActivity = useAtomValue(stationDashboardActivityAtom);

  return (
    <Box
      width={342}
      bgcolor={"#ffffff"}
      sx={{ border: "1px solid #00000014" }}
      p={2}
    >
      <Stack spacing={1}>
        <Stack direction={"row"} justifyContent={"space-between"}>
          <Typography variant="subtitle2">Vehicle Information</Typography>
          <Button size="small" variant="outlined">
            Report Issue
          </Button>
        </Stack>
        <Divider />
        {[
          {
            label: "License Plate",
            value: dashboardActivity?.vehicleInfo?.licensePlate,
          },
          {
            label: "Station Site",
            value: dashboardActivity?.vehicleInfo?.stationSiteName,
          },
          {
            label: "Status",
            value: dashboardActivity?.vehicleInfo?.status,
          },
          {
            label: "Arrival Time",
            value: fDateTime(dashboardActivity?.arrivalTime),
          },
          {
            label: "License Plate Type",
            value: dashboardActivity?.vehicleInfo?.licensePlateType,
          },
          {
            label: "Vehicle Type",
            value: dashboardActivity?.vehicleInfo?.vehicleType,
          },
          {
            label: "Brand",
            value: dashboardActivity?.vehicleInfo?.band,
          },
          {
            label: "Color",
            value: dashboardActivity?.vehicleInfo?.colorName,
            type: "COLOR",
          },
          {
            label: "Tag",
            value: dashboardActivity?.tags
              ?.map((tag) => `${tag.tag}`)
              .join(", "),
            type: "TAGS",
          },
        ].map((info, i) => (
          <Stack
            key={i}
            direction={"row"}
            justifyContent={"space-between"}
            alignItems={"center"}
          >
            <Typography variant="subtitle2" width={130}>
              {info.label}
            </Typography>
            <Box flex={1} sx={{ border: "1px solid #00000014" }} p={1}>
              <Typography variant="subtitle2" fontWeight={400}>
                {info.value}
              </Typography>
            </Box>
          </Stack>
        ))}
      </Stack>
    </Box>
  );
}
