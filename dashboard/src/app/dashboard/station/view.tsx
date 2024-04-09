"use client";
import { useBoolean } from "@/hooks/use-boolean";
import { Add } from "@mui/icons-material";
import { Stack, Button, Modal, Box } from "@mui/material";
import React from "react";
import StationFilter from "./filter";
import StationTable from "./table";
import StationForm from "./station-form";

export default function StationView() {
  const openStationForm = useBoolean();
  return (
    <Stack spacing={1} p={2}>
      <StationForm
        opened={openStationForm.value}
        onClose={openStationForm.onFalse}
      />
      <StationFilter
        actions={[
          // eslint-disable-next-line react/jsx-key
          <Button
            variant="contained"
            color="primary"
            onClick={openStationForm.onTrue}
          >
            <Add />
            Create
          </Button>,
        ]}
      />
      <Box sx={{ flex: 1 }}>
        <StationTable />
      </Box>
    </Stack>
  );
}
