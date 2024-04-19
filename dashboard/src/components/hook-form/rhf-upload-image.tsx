import { Controller, useFormContext } from "react-hook-form";

import Box from "@mui/material/Box";
import Chip from "@mui/material/Chip";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import Checkbox from "@mui/material/Checkbox";
import InputLabel from "@mui/material/InputLabel";
import { Theme, SxProps } from "@mui/material/styles";
import FormHelperText from "@mui/material/FormHelperText";
import TextField, { TextFieldProps } from "@mui/material/TextField";
import FormControl, { FormControlProps } from "@mui/material/FormControl";
import { Button, Stack, Typography, styled } from "@mui/material";
import Image from "@/components/image";
import { AddPhotoAlternateOutlined } from "@mui/icons-material";
import { ChangeEvent, ChangeEventHandler, useRef } from "react";
import { UploadResult, useUploadFileMutation } from "@/apollo-client";
// ----------------------------------------------------------------------

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
type RHFUploadImageProps = {
  previewImage?: string;
  onUploaded?: (result: UploadResult) => void;
};

export function RHFUploadImage({
  previewImage,
  onUploaded,
  ...other
}: RHFUploadImageProps) {
  const fileInputRef = useRef<HTMLLabelElement | null>(null);
  const [uploadfile] = useUploadFileMutation();

  const handleUploadImage = async (event?: ChangeEvent<HTMLInputElement>) => {
    if (event?.target.files?.length) {
      const file = event?.target?.files[0]; // Get the selected file
      if (file && file.type.startsWith("image/")) {
        const files = await uploadfile({
          variables: {
            file: file,
          },
        });
        if (files.data?.uploadFile) {
          onUploaded && onUploaded(files.data?.uploadFile);
        }
      }
    }
  };
  if (previewImage) {
    return (
      <Stack
        alignItems={"center"}
        justifyContent={"center"}
        sx={{
          bgcolor: "rgba(250, 249, 249, 1)",
          border: "1px solid rgba(0, 0, 0, 0.05)",
          stroke: "rgb(28 27 31 / 41%)",
        }}
        height={250}
      >
         <Button
          ref={fileInputRef}
          component="label"
          role={undefined}
          tabIndex={-1}
        >
        <Image
          src={previewImage}
          alt="Preview"
          sx={{
            width: 330,
            height: 250,
          }}
        />     <VisuallyHiddenInput type="file" onChange={handleUploadImage} />
        </Button>
      </Stack>
    );
  }

  return (
    <Box height={1}>
      <Stack
        alignItems={"center"}
        justifyContent={"center"}
        sx={{
          bgcolor: "rgba(250, 249, 249, 1)",
          border: "1px solid rgba(0, 0, 0, 0.05)",
          stroke: "rgb(28 27 31 / 41%)",
        }}
        height={1}
      >
        <AddPhotoAlternateOutlined
          style={{
            fontSize: "28px",
            marginBottom: "28px",
          }}
        />

        <Button
          ref={fileInputRef}
          component="label"
          role={undefined}
          // variant="contained"
          tabIndex={-1}
          // type="file"
          // onClick={handleButtonClick}
        >
          <Typography
            sx={{
              border: "1px solid rgba(0, 0, 0, 0.05)",
              borderRadius: "4px",
              padding: "4px",
              color: "rgba(0, 0, 0, 0.5)",
            }}
          >
            + Add image
            <VisuallyHiddenInput type="file" onChange={handleUploadImage} />
          </Typography>
        </Button>
      </Stack>
    </Box>
  );
}

// ----------------------------------------------------------------------

type RHFMultiSelectProps = FormControlProps & {
  name: string;
  label?: string;
  chip?: boolean;
  checkbox?: boolean;
  placeholder?: string;
  helperText?: React.ReactNode;
  options: {
    label: string;
    value: string;
  }[];
};

export function RHFMultiSelect({
  name,
  chip,
  label,
  options,
  checkbox,
  placeholder,
  helperText,
  ...other
}: RHFMultiSelectProps) {
  const { control } = useFormContext();

  const renderValues = (selectedIds: string[]) => {
    const selectedItems = options.filter((item) =>
      selectedIds.includes(item.value)
    );

    if (!selectedItems.length && placeholder) {
      return <Box sx={{ color: "text.disabled" }}>{placeholder}</Box>;
    }

    if (chip) {
      return (
        <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
          {selectedItems.map((item) => (
            <Chip key={item.value} size="small" label={item.label} />
          ))}
        </Box>
      );
    }

    return selectedItems.map((item) => item.label).join(", ");
  };

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState: { error } }) => (
        <FormControl error={!!error} {...other}>
          {label && <InputLabel id={name}> {label} </InputLabel>}

          <Select
            {...field}
            multiple
            displayEmpty={!!placeholder}
            id={`multiple-${name}`}
            labelId={name}
            label={label}
            renderValue={renderValues}
          >
            {options.map((option) => {
              const selected = field.value.includes(option.value);

              return (
                <MenuItem key={option.value} value={option.value}>
                  {checkbox && (
                    <Checkbox size="small" disableRipple checked={selected} />
                  )}

                  {option.label}
                </MenuItem>
              );
            })}
          </Select>

          {(!!error || helperText) && (
            <FormHelperText error={!!error}>
              {error ? error?.message : helperText}
            </FormHelperText>
          )}
        </FormControl>
      )}
    />
  );
}
