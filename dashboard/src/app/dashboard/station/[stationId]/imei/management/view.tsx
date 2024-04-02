"use client";
import { useBoolean } from "@/hooks/use-boolean";
import { Add } from "@mui/icons-material";
import { Stack, Button, Box } from "@mui/material";
import React from "react";
import ImeiFilter from "./filter";
import ImeiTable from "./table";
import ImeiForm from "./imei-form";
import { useGetMobileDeviceConfigurationsQuery } from "@/apollo-client";

type Props = {
  stationLocationId: string;
};
export default function ImeiManageView({ stationLocationId }: Props) {
  const openImeiManageForm = useBoolean();
  const { data, loading } = useGetMobileDeviceConfigurationsQuery({
    variables: {
      stationLocationId,
      limit: 100000,
      skip: 0,
    },
  });
  return (
    <Stack spacing={1}>
      <ImeiFilter
        actions={[
          // eslint-disable-next-line react/jsx-key
          <Button
            variant="contained"
            color="primary"
            onClick={openImeiManageForm.onTrue}
          >
            <Add />
            Create
          </Button>,
        ]}
      />
      <ImeiForm
        opened={openImeiManageForm.value}
        onClose={openImeiManageForm.onFalse}
        dafaultValues={{ stationLocationId }}
      />
      <Box sx={{ flex: 1 }}>
        <ImeiTable data={data} loading={loading} />
      </Box>
    </Stack>
  );
}
