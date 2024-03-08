"use client";

import { AuthGuard, GuestGuard } from "@/auth/guard";
// ----------------------------------------------------------------------

type Props = {
  children: React.ReactNode;
};

export default function Layout({ children }: Props) {
  return <GuestGuard>{children}</GuestGuard>;
}
