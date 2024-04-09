import { Controller, useFormContext } from "react-hook-form";

import Chip from "@mui/material/Chip";
import TextField from "@mui/material/TextField";
import InputAdornment from "@mui/material/InputAdornment";
import Autocomplete, { AutocompleteProps } from "@mui/material/Autocomplete";

import Iconify from "@/components/iconify";

// ----------------------------------------------------------------------

interface Props<
  T,
  Multiple extends boolean | undefined,
  DisableClearable extends boolean | undefined,
  FreeSolo extends boolean | undefined
> extends AutocompleteProps<T, Multiple, DisableClearable, FreeSolo> {
  name: string;
  label?: string;
  placeholder?: string;
  type?: "country" | string;
  helperText?: React.ReactNode;
}

export default function RHFAutocomplete<
  T,
  Multiple extends boolean | undefined,
  DisableClearable extends boolean | undefined,
  FreeSolo extends boolean | undefined
>({
  name,
  label,
  type,
  helperText,
  placeholder,
  ...other
}: Omit<Props<T, Multiple, DisableClearable, FreeSolo>, "renderInput">) {
  const { control, setValue } = useFormContext();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState: { error } }) => {
        return (
          <Autocomplete
            {...field}
            id={`autocomplete-${name}`}
            onChange={(event, newValue) =>
              setValue(name, newValue, { shouldValidate: true })
            }
            renderInput={(params) => (
              <TextField
                {...params}
                label={label}
                placeholder={placeholder}
                error={!!error}
                helperText={error ? error?.message : helperText}
                inputProps={{
                  ...params.inputProps,
                }}
              />
            )}
            {...other}
          />
        );
      }}
    />
  );
}

// ----------------------------------------------------------------------

export function getCountry(inputValue: string) {
  const option = [].filter(
    (country: any) => country.label === inputValue
  )[0] as any;

  return {
    ...option,
  };
}
