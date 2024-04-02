import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import TableRow from "@mui/material/TableRow";

import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";

import { useBoolean } from "@/hooks/use-boolean";
import Iconify from "@/components/iconify";

import { IVehicleItem } from "./types";
import CustomPopover, { usePopover } from "@/components/custom-popover";
import Label from "@/components/label";
import { ConfirmDialog } from "@/components/custom-dialog";
import { Delete, Edit } from "@mui/icons-material";
import { Alert, Typography } from "@mui/material";
import { useRouter } from "@/routes/hooks";
import { DevicePermittedLabel } from "@/apollo-client";

// ----------------------------------------------------------------------

type Props = {
  row: IVehicleItem;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

export default function VehicleTableRow({ row, onDeleteRow, onEdit }: Props) {
  const { number, prefix, tags, type, permittedLabel } = row;
  const confirm = useBoolean();
  return (
    <>
      <TableRow hover>
        <TableCell>
          <Typography fontSize={14} color={"text.primary"}>
            -
          </Typography>
        </TableCell>
        <TableCell>
          <Stack direction={"row"} spacing={0.5}>
            <Typography fontSize={14} color={"text.primary"}>
              {prefix}
            </Typography>
            <Typography fontSize={14} color={"text.primary"}>
              {number}
            </Typography>
          </Stack>
        </TableCell>
        <TableCell>
          <Typography fontSize={14} color={"text.primary"}>
            License Plate Type
          </Typography>
        </TableCell>
        <TableCell>
          <Typography fontSize={14} color={"text.primary"}>
            Band
          </Typography>
        </TableCell>
        <TableCell>
          <Typography fontSize={14} color={"text.primary"}>
            {type}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography fontSize={14} color={"text.primary"}>
            Color
          </Typography>
        </TableCell>
        <TableCell>
          <Label
            variant="filled"
            borderRadius={1}
            color={
              (permittedLabel === DevicePermittedLabel.Whitelist &&
                "success") ||
              (permittedLabel === DevicePermittedLabel.Blacklist && "error") ||
              "default"
            }
          >
            {permittedLabel}
          </Label>
        </TableCell>
        <TableCell>
          <Stack
            alignItems={"center"}
            spacing={1}
            direction={"row"}
            flexWrap={"wrap"}
          >
            {tags?.map((tag) => (
              <Label key={tag.id}>{tag.title} </Label>
            ))}
          </Stack>
        </TableCell>
        <TableCell>
          <Stack
            direction={"row"}
            alignItems={"center"}
            justifyContent={"center"}
          >
            <IconButton size="small" onClick={onEdit}>
              <Edit fontSize="small" />
            </IconButton>
            <IconButton size="small" onClick={confirm.onTrue}>
              <Delete fontSize="small" />
            </IconButton>
          </Stack>
        </TableCell>
      </TableRow>

      <ConfirmDialog
        open={confirm.value}
        onClose={confirm.onFalse}
        title="Delete"
        content="Are you sure want to delete?"
        action={
          <Button variant="contained" color="error" onClick={onDeleteRow}>
            Delete
          </Button>
        }
      />
    </>
  );
}
