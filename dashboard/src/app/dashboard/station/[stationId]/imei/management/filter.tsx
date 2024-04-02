import { Card, Typography, Divider, Box, Stack } from "@mui/material";
import React from "react";
type Props = {
  actions?: JSX.Element[];
};
export default function ImeiFilter({ actions }: Props) {
  return (
    <Card>
      <Box
        sx={{
          height: 56,
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
        px={2}
      >
        <Typography variant="h5">IMEI/IMSI</Typography>
        {actions && (
          <Stack alignItems={"center"} spacing={1}>
            {actions.map((action, i) => (
              <React.Fragment key={i}>{action}</React.Fragment>
            ))}
          </Stack>
        )}
      </Box>
      <Divider />
      <Box sx={{ height: 44, display: "flex", alignItems: "center" }} px={2}>
        {/* <Typography variant="body2">Station Site</Typography>
         */}
      </Box>
    </Card>
  );
}
