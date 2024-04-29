import { Stack, Typography, Button, Divider, Box } from "@mui/material";
import { useAtomValue } from "jotai";
import React from "react";
import { stationDashboardActivityAtom } from "./store";
import { fDateTime } from "@/utils/format-time";
import { StationDashboardActivity } from "@/apollo-client";
import CreateReportModal from "./report-modal";
import { useBoolean } from "@/hooks/use-boolean";

type Props = {
  vehicleInfo? : StationDashboardActivity
}

export default function VehicleInfo({vehicleInfo} : Props) {
  const dashboardActivity = useAtomValue(stationDashboardActivityAtom);

  const openReportForm = useBoolean();

  const openReportModal = async () =>{
    openReportForm.onTrue()
  }

  return (
    <Box
      width={342}
      bgcolor={"#ffffff"}
      sx={{ border: "1px solid #00000014" }}
      p={2}
    >
            <CreateReportModal
        // row={vehicleInfo}
        // stationId={stationLocationId}
        opened={openReportForm.value}
        onClose={openReportForm.onFalse}
   
      />
      
      <Stack spacing={1}>
        <Stack direction={"row"} justifyContent={"space-between"}>
          <Typography variant="subtitle2">Vehicle Information</Typography>
          {!vehicleInfo && <Button size="small" variant="outlined" 
          onClick={openReportModal}
          >
            Report Issue
          </Button>}
        </Stack>
        <Divider />
        {[
          {
            label: "License Plate",
            value: vehicleInfo?.vehicleInfo?.licensePlate ?? dashboardActivity?.vehicleInfo?.licensePlate,
          },
          {
            label: "Station Site",
            value: vehicleInfo?.vehicleInfo?.stationSiteName ?? dashboardActivity?.vehicleInfo?.stationSiteName,
          },
          {
            label: "Vehicle Status",
            value: vehicleInfo?.vehicleInfo?.status ?? dashboardActivity?.vehicleInfo?.status,
          },
          {
            label: "Arrival Time",
            value: fDateTime(vehicleInfo?.arrivalTime ?? dashboardActivity?.arrivalTime),
          },
          {
            label: "License Plate Type",
            value: vehicleInfo?.vehicleInfo?.licensePlateType ?? dashboardActivity?.vehicleInfo?.licensePlateType,
          },
          {
            label: "Vehicle Type",
            value: vehicleInfo?.vehicleInfo?.vehicleType ?? dashboardActivity?.vehicleInfo?.vehicleType,
          },
          {
            label: "Brand",
            value: vehicleInfo?.vehicleInfo?.band ?? dashboardActivity?.vehicleInfo?.band,
          },
          {
            label: "Color",
            value: vehicleInfo?.vehicleInfo?.colorName ?? dashboardActivity?.vehicleInfo?.colorName,
            type: "COLOR",
          },
          {
            label: "Tag",
            value: (vehicleInfo?.tags ?? dashboardActivity?.tags)
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
