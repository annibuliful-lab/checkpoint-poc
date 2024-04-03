import { DevicePermittedLabel } from "@/apollo-client";
import { SimCard } from "@mui/icons-material";
import { Card, Stack, Typography } from "@mui/material";
import React from "react";
type Props = {
  notification: DevicePermittedLabel;
};
export default function NotificationItem() {
  return (
    <Card sx={{ border: "1px solid #00000014", p: 1 }}>
      <Stack direction={"row"} spacing={2}>
        <Stack
          alignItems={"center"}
          justifyContent={"center"}
          height={32}
          sx={{ bgcolor: "error.dark" }}
          width={32}
          borderRadius={2}
        >
          <SimCard htmlColor="white" />
        </Stack>
        <Stack>
          <Typography fontWeight={600} variant="subtitle2">
            Received Danger IMSI: 520031234567890
          </Typography>
          <Typography fontWeight={400} variant="subtitle2">
            4กข 490 กรุงเทพมหานคร, HONDA
          </Typography>
          <Typography fontWeight={400} variant="subtitle2">
            Black, IMEI: 861536030196001
          </Typography>
          <Typography mt={1.5} variant="caption">
            12/02/2024 12:00:22
          </Typography>
        </Stack>
      </Stack>
    </Card>
  );
}
export function NotificationPermittedLabel({
  devicePermittedLabel,
}: {
  devicePermittedLabel: DevicePermittedLabel;
}) {
  return (
    <Card sx={{ border: "1px solid #00000014", p: 1 }}>
      <Stack direction={"row"} spacing={2}>
        <Stack
          alignItems={"center"}
          justifyContent={"center"}
          height={32}
          sx={{ bgcolor: "error.dark" }}
          width={32}
          borderRadius={2}
        >
          <SimCard htmlColor="white" />
        </Stack>
        <Stack>
          <Typography fontWeight={600} variant="subtitle2">
            Received Danger IMSI: 520031234567890
          </Typography>
          <Typography fontWeight={400} variant="subtitle2">
            4กข 490 กรุงเทพมหานคร, HONDA
          </Typography>
          <Typography fontWeight={400} variant="subtitle2">
            Black, IMEI: 861536030196001
          </Typography>
          <Typography mt={1.5} variant="caption">
            12/02/2024 12:00:22
          </Typography>
        </Stack>
      </Stack>
    </Card>
  );
}
