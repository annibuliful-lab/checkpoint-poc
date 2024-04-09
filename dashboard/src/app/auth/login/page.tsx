"use client";

import { Box } from "@mui/material";
import LoginForm from "./login-form";

export default function Page() {
  return (
    <Box
      sx={{
        width: 1,
        height: 1,
        justifyContent: "center",
        alignItems: "center",
        display: "flex",
      }}
    >
      <LoginForm />
    </Box>
  );
}
