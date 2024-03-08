"use client";

import { useContext, useEffect, useState } from "react";

import { AuthContext } from "../context/jwt/auth-context";
import { useAtomValue } from "jotai/react";
import { authAtom } from "../store";

// ----------------------------------------------------------------------

export const useAuthContext = () => {
  const context = useContext(AuthContext);

  if (!context)
    throw new Error("useAuthContext context must be use inside AuthProvider");

  return context;
};

export const useAuthAtom = () => {
  const authenticated = useAtomValue(authAtom);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    setLoading(false);
  }, []);
  return {
    loading,
    authenticated,
  };
};
