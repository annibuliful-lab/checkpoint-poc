import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import TableRow from "@mui/material/TableRow";

import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";

import { useBoolean } from "@/hooks/use-boolean";
import Iconify from "@/components/iconify";

import { ImeiItem } from "./types";
import CustomPopover, { usePopover } from "@/components/custom-popover";
import Label from "@/components/label";
import { ConfirmDialog } from "@/components/custom-dialog";
import { Edit } from "@mui/icons-material";
import { Typography } from "@mui/material";
import { useRouter } from "@/routes/hooks";
import { DevicePermittedLabel } from "@/apollo-client";

// ----------------------------------------------------------------------

type Props = {
  row: ImeiItem;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

export default function ImeiTableRow({ row, onDeleteRow, onEdit }: Props) {
  const {
    blacklistPriority,
    createdAt,
    createdBy,
    id,
    msisdn,
    permittedLabel,
    projectId,
    referenceImeiConfiguration,
    referenceImeiConfigurationId,
    referenceImsiConfiguration,
    referenceImsiConfigurationId,
    stationLocationId,
    tags,
    title,
    updatedAt,
    updatedBy,
  } = row;
  const confirm = useBoolean();
  const popover = usePopover();
  const route = useRouter();
  return (
    <>
      <TableRow hover>
        <TableCell>
          <Typography color={"text.primary"}>
            {referenceImeiConfiguration?.imei}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography color={"text.primary"}>
            {referenceImsiConfiguration?.imsi}
          </Typography>
        </TableCell>
        <TableCell>
          <Label
            variant="soft"
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
          <Stack alignItems={"center"} justifyContent={"center"}>
            <IconButton size="small" onClick={onEdit}>
              <Edit fontSize="small" />
            </IconButton>
          </Stack>
        </TableCell>
      </TableRow>
      <CustomPopover
        open={popover.open}
        onClose={popover.onClose}
        arrow="right-top"
        sx={{ width: 140 }}
      >
        <MenuItem
          onClick={() => {
            confirm.onTrue();
            popover.onClose();
          }}
          sx={{ color: "error.main" }}
        >
          <Iconify icon="solar:trash-bin-trash-bold" />
          Delete
        </MenuItem>
      </CustomPopover>

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
