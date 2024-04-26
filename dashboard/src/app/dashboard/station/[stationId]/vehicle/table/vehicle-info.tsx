import {
  Box,
  Dialog, 
  DialogContent,
  DialogTitle,
  Stack,
} from "@mui/material";
import React from "react"; 
import _ from "lodash";
import VehicleCamera from "../../dashboard/vehicle-camera";
import VehicleInfo from "../../dashboard/vehicle-info";
import Transection from "../../dashboard/transaction";
import { VehicleTransection } from "./types";
type Props = {
  opened: boolean;
  onClose: () => void;
  stationId: string;
  row?: VehicleTransection;
};
export default function VehicleInfoModal({
  opened,
  onClose,
  stationId,
  row
}: Props) {
  const title = "Vehicle Information";
console.log({row})
  // IStationItem

  return (
    <Dialog
      fullWidth
      open={opened}
      onClose={onClose}
      maxWidth='xl'
    
    >
      <DialogTitle sx={{ pb: 2 }}>{title}</DialogTitle>
      <DialogContent
        sx={{
          overflow: "unset",
        }}
      >
        <Stack spacing={0.5}>
          <Stack direction={"row"} spacing={0.5}>
            <Box flex={1}>
              <VehicleCamera />
            </Box>
            <VehicleInfo 
            // VehicleInfo={

            // }
            />
          </Stack>
          <Transection stationLocationId={stationId} />
        </Stack>
      </DialogContent>
    </Dialog>
  );
}
