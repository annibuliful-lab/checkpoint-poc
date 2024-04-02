"use client";

import { Box, Stack } from "@mui/material";
import Transection from "./transaction";
import VehicleInfo from "./vehicle-info";
import VehicleCamera from "./vehicle-camera";

export default function VehicleView() {
  return (
    <Stack spacing={0.5}>
      <Stack direction={"row"} spacing={0.5}>
        <Box flex={1}>
          <VehicleCamera />
        </Box>
        <VehicleInfo />
      </Stack>
      <Transection />
    </Stack>
  );
}
