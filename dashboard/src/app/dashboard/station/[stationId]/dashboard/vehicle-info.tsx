import { Stack, Typography, Button, Divider, Box } from "@mui/material";
import React from "react";

export default function VehicleInfo() {
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
            value: "1กฮ 1234 กรุงเทพมหานคร",
          },
          {
            label: "Station Site",
            value: "ภ.9 ยะลา",
          },
          {
            label: "Status",
            value: "BLACKLIST",
          },
          {
            label: "Arrival Time",
            value: "12/02/2024 12:00:22",
          },
          {
            label: "License Plate Type",
            value: "Private Car",
          },
          {
            label: "Vehicle Type",
            value: "Standard Sedan",
          },
          {
            label: "Brand",
            value: "Honda",
          },
          {
            label: "Color",
            value: "Black",
          },
          {
            label: "Tag",
            value: "-",
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
