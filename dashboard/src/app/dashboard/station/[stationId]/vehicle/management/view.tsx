"use client";
import { useBoolean } from "@/hooks/use-boolean";
import { Add } from "@mui/icons-material";
import { Stack, Button, Box } from "@mui/material";
import React from "react";
import VehicleFilter from "./filter";
import VehicleTable from "./table";
import VehicleForm from "./vehicle-form";
import { PropWithStationLocationId } from "../../types";

export default function VehicleManageView({
  stationLocationId,
}: PropWithStationLocationId) {
  const openVehicleManageForm = useBoolean();
  return (
    <Stack spacing={1}>
      <VehicleFilter
        actions={[
          // eslint-disable-next-line react/jsx-key
          <Button
            variant="contained"
            color="primary"
            onClick={openVehicleManageForm.onTrue}
          >
            <Add />
            Create
          </Button>,
        ]}
      />
      <VehicleForm
        opened={openVehicleManageForm.value}
        onClose={openVehicleManageForm.onFalse}
        vehicle={{ stationLocationId: stationLocationId}}
      />
      <Box sx={{ flex: 1 }}>
        <VehicleTable />
      </Box>
    </Stack>
  );
}
