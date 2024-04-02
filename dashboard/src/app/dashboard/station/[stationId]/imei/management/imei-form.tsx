import FormProvider from "@/components/hook-form/form-provider";
import {
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  MenuItem,
  Stack,
  Typography,
} from "@mui/material";
import React, { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  RHFAutocomplete,
  RHFSelect,
  RHFTextField,
} from "@/components/hook-form";
import _ from "lodash";
import LoadingButton from "@mui/lab/LoadingButton";
import AlertApolloError from "@/components/alert-apollo-error";
import {
  ImeiFormInput,
  ImeiFormSchema,
  imeiFormDefaultValues,
} from "./imei-form.schema";
import {
  MobileDeviceConfiguration,
  useCreateMobileDeviceConfigurationMutation,
  useGetTagsQuery,
  useUpdateMobileDeviceConfigurationMutation,
} from "@/apollo-client";
import { IMEI_PRIORITY_OPTIONS, IMEI_STATUS_OPTIONS } from "./const";
type Props = {
  opened: boolean;
  onClose: () => void;
  dafaultValues?: Partial<MobileDeviceConfiguration>;
};
export default function ImeiForm({ opened, onClose, dafaultValues }: Props) {
  const [create, createResponse] = useCreateMobileDeviceConfigurationMutation();
  const [update, updateResponse] = useUpdateMobileDeviceConfigurationMutation();
  const { data } = useGetTagsQuery({
    variables: {
      limit: 100,
      skip: 0,
    },
  });
  const tags = useMemo(
    () => data?.getTags.map((t) => t.title) ?? [],
    [data?.getTags]
  );
  const initialVal = useMemo(
    (): ImeiFormInput => ({
      imei:
        dafaultValues?.referenceImeiConfiguration?.imei ??
        imeiFormDefaultValues.imei,
      imsi:
        dafaultValues?.referenceImsiConfiguration?.imsi ??
        imeiFormDefaultValues.imsi,
      tags:
        dafaultValues?.tags?.map((tag) => tag.title) ??
        imeiFormDefaultValues.tags,
      stationLocationId:
        dafaultValues?.stationLocationId ??
        imeiFormDefaultValues.stationLocationId,
      title: dafaultValues?.title ?? imeiFormDefaultValues.title,
      permittedLabel:
        dafaultValues?.permittedLabel ?? imeiFormDefaultValues.permittedLabel,
      blacklistPriority:
        dafaultValues?.blacklistPriority ??
        imeiFormDefaultValues.blacklistPriority,
    }),
    [
      dafaultValues?.blacklistPriority,
      dafaultValues?.permittedLabel,
      dafaultValues?.referenceImeiConfiguration?.imei,
      dafaultValues?.referenceImsiConfiguration?.imsi,
      dafaultValues?.stationLocationId,
      dafaultValues?.tags,
      dafaultValues?.title,
    ]
  );
  const methods = useForm({
    resolver: yupResolver(ImeiFormSchema),
    defaultValues: initialVal,
  });

  const {
    reset,
    handleSubmit,
    formState: { isSubmitting },
  } = methods;
  const onSubmit = handleSubmit(async (values) => {
    try {
      const input = values;
      if (dafaultValues?.id) {
        await update({ variables: { id: dafaultValues?.id, ...input } });
        reset(input);
      } else {
        await create({ variables: input });
        reset();
      }
      onClose();
    } catch (error) {
      console.error(error);
    }
    createResponse.reset();
    updateResponse.reset();
  });

  useEffect(() => {
    if (dafaultValues) {
      reset(initialVal);
    }

    return () => {};
  }, [dafaultValues, initialVal, reset]);

  const title = "Imei infomation";

  return (
    <Dialog
      fullWidth
      open={opened}
      onClose={onClose}
      PaperProps={{
        sx: {
          maxWidth: 544,
        },
      }}
    >
      <FormProvider methods={methods} onSubmit={onSubmit}>
        <DialogTitle sx={{ pb: 2 }}>{title}</DialogTitle>

        <DialogContent
          sx={{
            overflow: "unset",
          }}
        >
          <Stack spacing={2}>
            <AlertApolloError
              error={createResponse.error || updateResponse.error}
            />
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                IMEI Number
              </Typography>
              <RHFTextField
                sx={{ flex: 1 }}
                size="small"
                name="imei"
                placeholder="IMEI Number"
              />
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                IMSI Number
              </Typography>
              <RHFTextField
                sx={{ flex: 1 }}
                size="small"
                name="imsi"
                placeholder="IMSI Number"
              />
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                Status
              </Typography>
              <Stack
                direction={"row"}
                alignItems={"center"}
                spacing={2}
                sx={{ flex: 1 }}
              >
                <RHFSelect name="permittedLabel" label="Status" size="small">
                  {IMEI_STATUS_OPTIONS.map((option) => (
                    <MenuItem key={option.id} value={option.id}>
                      {option.label}
                    </MenuItem>
                  ))}
                </RHFSelect>
                <RHFSelect
                  name="blacklistPriority"
                  label="Priority"
                  size="small"
                >
                  {IMEI_PRIORITY_OPTIONS.map((option) => (
                    <MenuItem key={option.id} value={option.id}>
                      {option.label}
                    </MenuItem>
                  ))}
                </RHFSelect>
              </Stack>
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                Tags
              </Typography>
              <RHFAutocomplete
                freeSolo
                multiple
                fullWidth
                size="small"
                name="tags"
                placeholder="Tags"
                options={tags}
                sx={{ flex: 1 }}
                renderTags={(tagValue, getTagProps) => {
                  return tagValue.map((option, index) => (
                    <Chip
                      {...getTagProps({ index })}
                      key={option}
                      label={option}
                    />
                  ));
                }}
              />
            </Stack>
            <Divider flexItem sx={{ borderStyle: "dashed" }} />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button variant="outlined" onClick={onClose}>
            Cancel
          </Button>

          <LoadingButton
            type="submit"
            variant="contained"
            color="primary"
            loading={isSubmitting}
          >
            Submit
          </LoadingButton>
        </DialogActions>
      </FormProvider>
    </Dialog>
  );
}
