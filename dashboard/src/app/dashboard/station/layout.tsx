"use client";
import { Stack } from "@mui/material";
// ----------------------------------------------------------------------

type Props = {
  children: React.ReactNode;
};

export default function Layout({ children }: Props) {
  return <Stack>{children}</Stack>;
}
