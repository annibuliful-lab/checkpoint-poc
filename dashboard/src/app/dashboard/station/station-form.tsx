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
  IconButton,
  InputAdornment,
  Stack,
  Typography,
} from "@mui/material";
import React, { useEffect, useMemo } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import {
  StationFormInput,
  StationFormSchema,
  stationFormDefaultValues,
} from "./station-form.schema";
import { RHFAutocomplete, RHFTextField } from "@/components/hook-form";
import {
  GetStationLocationsDocument,
  StationLocation,
  Tag,
  useCreateStationLocationMutation,
  useGetTagsQuery,
  useUpdateStationLocationMutation,
} from "@/apollo-client";
import _ from "lodash";
import LoadingButton from "@mui/lab/LoadingButton";
import { AddLocationAlt } from "@mui/icons-material";
import { useBoolean } from "@/hooks/use-boolean";
import Iconify from "@/components/iconify";
import AlertApolloError from "@/components/alert-apollo-error";
type Props = {
  opened: boolean;
  onClose: () => void;
  station?: StationLocation;
};
export default function StationForm({ opened, onClose, station }: Props) {
  const [create, createResponse] = useCreateStationLocationMutation();
  const [update, updateResponse] = useUpdateStationLocationMutation();
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
  const onAddLocation = useBoolean();
  const initialVal = useMemo(
    () =>
      ({
        title: station?.title ?? stationFormDefaultValues.title,
        description:
          station?.description ?? stationFormDefaultValues.description,
        department: station?.department ?? stationFormDefaultValues.department,
        latitude: station?.latitude ?? stationFormDefaultValues.latitude,
        longitude: station?.longitude ?? stationFormDefaultValues.longitude,
        tags:
          station?.tags?.map((tag) => tag.title) ??
          stationFormDefaultValues.tags,
      } as StationFormInput),
    [
      station?.department,
      station?.description,
      station?.latitude,
      station?.longitude,
      station?.tags,
      station?.title,
    ]
  );
  const methods = useForm({
    resolver: yupResolver(StationFormSchema),
    defaultValues: initialVal,
  });

  const {
    reset,
    control,
    handleSubmit,
    setValue,
    formState: { isSubmitting },
  } = methods;

  const { fields, append, remove } = useFieldArray({
    control,
    name: "responsiblePersons",
  });

  const handleAdd = () => {
    append({
      name: "",
      phone: "",
    });
  };
  const handleRemove = (index: number) => {
    remove(index);
  };

  const onSubmit = handleSubmit(async (values) => {
    const input = values as StationFormInput;
    if (station?.id) {
      await update({
        variables: { ...input, id: station.id },
        refetchQueries: [GetStationLocationsDocument],
      });
      reset(input);
    } else {
      await create({
        variables: input,
        refetchQueries: [GetStationLocationsDocument],
      });
      reset();
    }
    onClose();
  });

  useEffect(() => {
    if (station) {
      reset(initialVal);
    }

    return () => {};
  }, [station, reset, initialVal]);

  const title = "Station infomation";

  return (
    <Dialog
      fullWidth
      open={opened}
      onClose={onClose}
      PaperProps={{
        sx: {
          maxWidth: 544,
          minHeight: 522,
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
                Title
              </Typography>
              <RHFTextField size="small" name="title" placeholder="Title" />
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                Description
              </Typography>
              <RHFTextField
                size="small"
                name="description"
                placeholder="Description"
              />
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                Department
              </Typography>
              <RHFTextField
                size="small"
                name="department"
                placeholder="Department"
              />
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                Geolocation
              </Typography>
              <RHFTextField
                size="small"
                name="geolocation"
                disabled
                placeholder="Geolocation"
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        size="small"
                        onClick={onAddLocation.onTrue}
                        edge="end"
                      >
                        <AddLocationAlt fontSize="small" />
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />
            </Stack>
            <Stack direction={"row"} alignItems={"center"}>
              <Typography width={120} variant="subtitle2">
                Remark
              </Typography>
              <RHFTextField size="small" name="remark" placeholder="Remark" />
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
            <Stack
              divider={<Divider flexItem sx={{ borderStyle: "dashed" }} />}
              spacing={1.5}
            >
              <Stack
                direction={"row"}
                justifyContent={"space-between"}
                alignItems={"center"}
              >
                <Typography variant="subtitle2">Responsible Persons</Typography>
                <IconButton size="small" color="default" onClick={handleAdd}>
                  <Iconify icon="mingcute:add-line" />
                </IconButton>
              </Stack>
              {fields.map((item, index) => (
                <Stack key={item.id} alignItems="flex-end">
                  <Stack
                    direction={{ xs: "column", md: "row" }}
                    spacing={1}
                    sx={{ width: 1 }}
                  >
                    <RHFTextField
                      size="small"
                      name={`responsiblePersons[${index}].name`}
                      label="Name"
                      InputLabelProps={{ shrink: true }}
                    />
                    <RHFTextField
                      size="small"
                      name={`responsiblePersons[${index}].phone`}
                      label="Phone"
                      InputLabelProps={{ shrink: true }}
                    />
                    <IconButton
                      color="error"
                      onClick={() => handleRemove(index)}
                    >
                      <Iconify icon="solar:trash-bin-trash-bold" />
                    </IconButton>
                  </Stack>
                </Stack>
              ))}
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
