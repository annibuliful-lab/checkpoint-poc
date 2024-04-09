"use client";
import { useSigninMutation } from "@/apollo-client";
import { authAtom } from "@/auth/store";
import { RHFTextField } from "@/components/hook-form";
import FormProvider from "@/components/hook-form/form-provider";
import { yupResolver } from "@hookform/resolvers/yup";
import { LockOpen, Person } from "@mui/icons-material";
import LoadingButton from "@mui/lab/LoadingButton";
import { Box, Button, Card, Stack, Typography } from "@mui/material";
import { useSetAtom } from "jotai/react";
import React, { useCallback } from "react";
import { useForm } from "react-hook-form";

//---------------------------------------------------------------------

export default function LoginForm() {
  const [callLogin, callLoginResponse] = useSigninMutation();
  const setLogin = useSetAtom(authAtom);
  const methods = useForm({
    // resolver: yupResolver(StationFormSchema),
    defaultValues: {
      username: "userA1234",
      password: "12345678",
    },
  });

  const {
    reset,
    control,
    handleSubmit,
    setValue,
    formState: { isSubmitting },
  } = methods;

  const onSubmit = handleSubmit(async (values) => {
    const { data } = await callLogin({
      variables: values,
    });
    if (data?.signin) {
      setLogin({
        ...data?.signin,
        projectId: "246bb085-8ccc-4def-ac78-dc2ad5c7760b",
      });
    }
  });

  return (
    <Card
      sx={{
        width: 412,
        height: 316,
        boxShadow: "0px 20px 40px -7px #00000026",
        p: 3,
      }}
    >
      <FormProvider methods={methods} onSubmit={onSubmit}>
        <Stack spacing={3}>
          <Typography
            color={"primary"}
            fontSize={36}
            fontWeight={600}
            letterSpacing={3}
            textAlign={"center"}
          >
            CHECKPOINT
          </Typography>

          <Stack direction={"row"} spacing={1} alignItems={"center"}>
            <Person color="disabled" />
            <RHFTextField size="small" name="username" placeholder="Username" />
          </Stack>
          <Stack direction={"row"} spacing={1} alignItems={"center"}>
            <LockOpen color="disabled" />
            <RHFTextField
              type="password"
              size="small"
              name="password"
              placeholder="Password"
            />
          </Stack>

          <LoadingButton
            fullWidth
            loading={callLoginResponse.loading}
            type="submit"
            variant="contained"
            color="primary"
          >
            Login
          </LoadingButton>
        </Stack>
      </FormProvider>
    </Card>
  );
}
