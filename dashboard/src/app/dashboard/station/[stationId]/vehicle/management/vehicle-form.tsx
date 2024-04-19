/* eslint-disable react/jsx-no-undef */
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
  Grid,
  MenuItem,
  Stack,
  Typography,
  styled,
} from "@mui/material";
import React, { useEffect, useMemo, useState } from "react";
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
  ImageType,
  UploadResult,
  UpsertImageS3KeyInput,
  VehicleTargetConfiguration,
  useCreateVehicleTargetConfigurationMutation,
  useGetTagsQuery,
  useUpdateVehicleTargetConfigurationMutation,
  useUploadFileMutation,
} from "@/apollo-client";
import { AddPhotoAlternateOutlined } from "@mui/icons-material";
import Image from "@/components/image";
import { RHFUploadImage } from "@/components/hook-form/rhf-upload-image";

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

  const [uploadfile] = useUploadFileMutation();
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
    (): VehicleFormInput => ({
      stationLocationId:
        vehicle?.stationLocationId ??
        vehicleFormDefaultValues.stationLocationId,
      prefix: vehicle?.prefix ?? vehicleFormDefaultValues.prefix,
      number: vehicle?.number ?? vehicleFormDefaultValues.number,
      province: vehicle?.province ?? vehicleFormDefaultValues.province,
      country: vehicle?.country ?? vehicleFormDefaultValues.country,
      permittedLabel:
        vehicle?.permittedLabel ?? vehicleFormDefaultValues.permittedLabel,
      blacklistPriority:
        vehicle?.blacklistPriority ??
        vehicleFormDefaultValues.blacklistPriority,
      type: vehicle?.type ?? vehicleFormDefaultValues.type,
      tags: vehicle?.tags?.map((t) => t.title) ?? vehicleFormDefaultValues.tags,
      imageS3Keys: vehicle?.images?.length
        ? vehicle?.images.map((img) => ({
            id: img.id,
            type: img.type,
            s3Key: img.s3Key,
          }))
        : vehicleFormDefaultValues?.imageS3Keys ?? [],
      brand: vehicle?.brand ?? vehicleFormDefaultValues.brand,
      color: vehicle?.color ?? vehicleFormDefaultValues.color,
    }),
    [
      vehicle?.stationLocationId,
      vehicle?.prefix,
      vehicle?.number,
      vehicle?.province,
      vehicle?.country,
      vehicle?.permittedLabel,
      vehicle?.blacklistPriority,
      vehicle?.type,
      vehicle?.tags,
      vehicle?.images,
      vehicle?.brand,
      vehicle?.color,
    ]
  );
  console.log({ initialVal });
  const methods = useForm({
    resolver: yupResolver(VehicleFormSchema as any),
    defaultValues: { ...initialVal },
  });
  const {
    reset,
    handleSubmit,
    setValue,
    watch,
    formState: { isSubmitting },
  } = methods;

  const [imagePreview, setImagePreview] = useState<
    Record<string, UploadResult>
  >({});
  const setUploadImageHandle = (result: UploadResult, i: number) => {
    const tempImagePreview : any[] = [];
    if (vehicle?.images?.length)
      for (let index = 0; index < vehicle?.images.length; index++) {
        const img = vehicle?.images?.[index];
        tempImagePreview[index] = {
          id: img.id || "",
          s3Key: img.s3Key,
          url: img.url,
        };
      }

    tempImagePreview[i] = {
      id: tempImagePreview?.[i]?.id || "",
      s3Key: result.s3Key,
      url: result.url,
    };
    setImagePreview((prevs) => ({ ...prevs, [i]: tempImagePreview[i] }));
    if (tempImagePreview?.[i]?.id) {
      setValue(`imageS3Keys.${i}.id`, tempImagePreview?.[i]?.id);
    }
    setValue(`imageS3Keys.${i}.s3Key`, tempImagePreview?.[i]?.s3Key);
    setValue(`imageS3Keys.${i}.type`, ImageType.Front);
  };

  const onSubmit = handleSubmit(async (values: VehicleFormInput) => {
    values.imageS3Keys = (values.imageS3Keys as UpsertImageS3KeyInput[]).filter(
      (img) => !!img.s3Key
    );
    const input = values;

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

  const handleOnClose = () => {
    reset();
    onClose();
  };

  const VisuallyHiddenInput = styled("input")({
    clip: "rect(0 0 0 0)",
    clipPath: "inset(50%)",
    height: 1,
    overflow: "hidden",
    position: "absolute",
    bottom: 0,
    left: 0,
    whiteSpace: "nowrap",
    width: 1,
  });

  useEffect(() => {
    if (vehicle) {
      reset(initialVal);
    }
    setImagePreview({});
    return () => {};
  }, [vehicle, reset, initialVal]);

  // useEffect(() => {
  //   if (vehicle?.images?.length) {
  //     for (const img of vehicle.images ?? []) {
  //       setImagePreview((prevs) => ({
  //         ...prevs,
  //         [img.s3Key]: { id: img.id || "", s3Key: img.s3Key, url: img.url },
  //       }));
  //     }
  //   }
  // }, [vehicle?.images]);

  const title = "Vehicle infomation";

  return (
    <Dialog fullWidth open={opened} onClose={onClose} maxWidth={"xl"}>
      <FormProvider methods={methods} onSubmit={onSubmit}>
        <DialogTitle key={1} sx={{ pb: 2 }}>
          {title}
        </DialogTitle>

        <DialogContent
          sx={{
            overflow: "unset",
          }}
        >
          <Stack spacing={2}>
            <Grid container spacing={1}>
              <Grid item sm={12} lg={8}>
                <Grid
                  container
                  height={1}
                  sx={{ border: "1px solid #0000001F" }}
                >
                  <Grid container height={1} spacing={1} p={1}>
                    {[...new Array(4)].map((_, i) => (
                      <Grid sm={12} lg={6} item key={i}>
                        <RHFUploadImage
                          previewImage={
                            imagePreview?.[i]?.url ||
                            vehicle?.images?.[i]?.url ||
                            ""
                          }
                          onUploaded={(result) => {
                            setUploadImageHandle(result, i);
                          }}
                        />
                      </Grid>
                    ))}
                  </Grid>
                </Grid>
              </Grid>
              <Grid item sm={12} lg={4}>
                <Stack spacing={2}>
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
                      <RHFTextField
                        size="small"
                        name="prefix"
                        label={"Prefix"}
                      />
                      <RHFTextField
                        size="small"
                        name="number"
                        label={"Number"}
                      />
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
                      <RHFTextField size="small" name="brand" label={"Brand"} />
                      <RHFTextField
                        size="small"
                        name="color"
                        label={"Color Name"}
                      />
                      {/* <RHFTextField
                        size="small"
                        name="vehicleType"
                        label={"Type"}
                      /> */}
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
                          {JSON.stringify(DevicePermittedLabel)}
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
              </Grid>
            </Grid>

            <Divider flexItem sx={{ borderStyle: "dashed" }} />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button variant="outlined" onClick={handleOnClose}>
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
