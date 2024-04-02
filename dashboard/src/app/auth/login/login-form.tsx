"use client";
import { useSigninMutation } from "@/apollo-client";
import { authAtom } from "@/auth/store";
import LoadingButton from "@mui/lab/LoadingButton";
import { Button } from "@mui/material";
import { useSetAtom } from "jotai/react";
import React, { useCallback } from "react";

//---------------------------------------------------------------------

export default function LoginForm() {
  const [callLogin, callLoginResponse] = useSigninMutation();
  const setLogin = useSetAtom(authAtom);
  const login = useCallback(async () => {
    const { data } = await callLogin({
      variables: {
        username: "userA1234",
        password: "12345678",
      },
    });
    if (data?.signin) {
      setLogin({
        ...data?.signin,
        projectId: "246bb085-8ccc-4def-ac78-dc2ad5c7760b",
      });
    }
  }, [callLogin, setLogin]);

  return (
    <div>
      <LoadingButton
        loading={callLoginResponse.loading}
        onClick={login}
        variant="contained"
      >
        Login
      </LoadingButton>
    </div>
  );
}
