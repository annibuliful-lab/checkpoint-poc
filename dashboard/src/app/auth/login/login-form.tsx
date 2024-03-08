"use client";
import { authAtom } from "@/auth/store";
import { Button } from "@mui/material";
import { useSetAtom } from "jotai/react";
import React, { useCallback } from "react";

//---------------------------------------------------------------------

export default function LoginForm() {
  // const [callLogin, callLoginResponse] = useSigninMutation();
  const setLogin = useSetAtom(authAtom);
  const login = useCallback(async () => {
    setLogin({
      refreshToken: "xxx",
      token: "xxx",
      userId: "xxx",
    });
  }, [setLogin]);

  return (
    <div>
      <Button onClick={login}>Login</Button>
    </div>
  );
}
