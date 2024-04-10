import FormProvider from "@/components/hook-form/form-provider";
import {
  Box,
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
  VehicleFormInput,
  VehicleFormSchema,
  vehicleFormDefaultValues,
} from "./vehicle-form.schema";
import {
  BlacklistPriority,
  DevicePermittedLabel,
  GetVehicleTargetConfigurationsDocument,
  VehicleTargetConfiguration,
  useCreateVehicleTargetConfigurationMutation,
  useGetTagsQuery,
  useUpdateVehicleTargetConfigurationMutation,
} from "@/apollo-client";

type Props = {
  opened: boolean;
  onClose: () => void;
  vehicle?: VehicleTargetConfiguration;
};
export default function VehicleForm({ opened, onClose, vehicle }: Props) {
  const [create, createResponse] =
    useCreateVehicleTargetConfigurationMutation();
  const [update, updateResponse] =
    useUpdateVehicleTargetConfigurationMutation();
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
    () => ({
      prefix: vehicle?.prefix ?? vehicleFormDefaultValues.prefix,
      number: vehicle?.number ?? vehicleFormDefaultValues.number,
      province: vehicle?.province ?? vehicleFormDefaultValues.province,
      type: vehicle?.type ?? vehicleFormDefaultValues.type,
      country: vehicle?.country ?? vehicleFormDefaultValues.country,
      permittedLabel:
        vehicle?.permittedLabel ?? vehicleFormDefaultValues.permittedLabel,
      blacklistPriority:
        vehicle?.blacklistPriority ??
        vehicleFormDefaultValues.blacklistPriority,
      tags: vehicle?.tags ?? vehicleFormDefaultValues.tags,
    }),
    [
      vehicle?.blacklistPriority,
      vehicle?.country,
      vehicle?.number,
      vehicle?.permittedLabel,
      vehicle?.prefix,
      vehicle?.province,
      vehicle?.tags,
      vehicle?.type,
    ]
  );
  const methods = useForm({
    resolver: yupResolver(VehicleFormSchema as any),
    defaultValues: initialVal,
  });

  const {
    reset,
    handleSubmit,
    formState: { isSubmitting },
  } = methods;
  const onSubmit = handleSubmit(async (values) => {
    const input = values as VehicleFormInput;
    if (vehicle?.id) {
      await update({
        variables: { ...input, id: vehicle.id },
        refetchQueries: [GetVehicleTargetConfigurationsDocument],
      });
      reset(input);
    } else {
      await create({
        variables: input,
        refetchQueries: [GetVehicleTargetConfigurationsDocument],
      });
      reset();
    }
    onClose();
  });

  useEffect(() => {
    if (vehicle) {
      reset(initialVal);
    }

    return () => {};
  }, [vehicle, initialVal, reset]);

  const title = "Vehicle infomation";

  return (
    <Dialog
      fullWidth
      open={opened}
      onClose={onClose}
      PaperProps={{
        sx: {
          maxWidth: 1232,
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
            <Stack direction={"row"} spacing={2}>
              <Box sx={{ width: 1 }}>Image</Box>
              <Stack spacing={2} sx={{ width: 1 }}>
                <AlertApolloError error={createResponse.error} />
                <Stack
                  direction={"row"}
                  alignItems={"center"}
                  height={60}
                  px={2}
                  sx={{ border: "1px solid #0000001F" }}
                >
                  <Typography width={120}>License Plate</Typography>
                  <Stack direction={"row"} spacing={1} flex={1}>
                    <RHFTextField size="small" name="prefix" label={"Prefix"} />
                    <RHFTextField size="small" name="number" label={"Number"} />
                    <RHFTextField
                      size="small"
                      name="province"
                      label={"Province"}
                    />
                    <RHFTextField size="small" name="type" label={"Type"} />
                  </Stack>
                </Stack>

                <Stack
                  direction={"row"}
                  alignItems={"center"}
                  height={60}
                  px={2}
                  sx={{ border: "1px solid #0000001F" }}
                >
                  <Typography width={120}>Vehicle Info</Typography>
                  <Stack direction={"row"} spacing={1} flex={1}>
                    <RHFTextField
                      size="small"
                      name="vehicleBrand"
                      label={"Brand"}
                    />
                    <RHFTextField
                      size="small"
                      name="vehicleColorName"
                      label={"Color Name"}
                    />
                    <RHFTextField
                      size="small"
                      name="vehicleType"
                      label={"Type"}
                    />
                  </Stack>
                </Stack>
                <Stack height={300} sx={{ border: "1px solid #0000001F" }}>
                  <Stack
                    direction={"row"}
                    alignItems={"center"}
                    height={60}
                    px={2}
                  >
                    <Typography width={120}>Status</Typography>
                    <Stack direction={"row"} spacing={1} flex={1}>
                      <RHFSelect
                        name="permittedLabel"
                        label="Status"
                        size="small"
                      >
                        {Object.keys(DevicePermittedLabel).map((key) => (
                          <MenuItem
                            key={key}
                            value={DevicePermittedLabel[key as "None"]}
                          >
                            {key}
                          </MenuItem>
                        ))}
                      </RHFSelect>
                      <RHFSelect
                        name="blacklistPriority"
                        label="Priority"
                        size="small"
                      >
                        {Object.keys(BlacklistPriority).map((key) => (
                          <MenuItem
                            key={key}
                            value={BlacklistPriority[key as "Normal"]}
                          >
                            {key}
                          </MenuItem>
                        ))}
                      </RHFSelect>
                    </Stack>
                  </Stack>
                  <Stack
                    direction={"row"}
                    alignItems={"center"}
                    height={60}
                    px={2}
                  >
                    <Typography width={120}>Tags</Typography>
                    <Stack direction={"row"} spacing={1} flex={1}>
                      <RHFAutocomplete
                        freeSolo
                        multiple
                        fullWidth
                        size="small"
                        name="tags"
                        placeholder="Tags"
                        options={tags}
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
                  </Stack>
                </Stack>
              </Stack>
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
