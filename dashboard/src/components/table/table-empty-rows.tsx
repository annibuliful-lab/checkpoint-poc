import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";

// ----------------------------------------------------------------------

type Props = {
  height?: number;
  emptyRows: number;
  colSpan?: number;
};

export default function TableEmptyRows({
  emptyRows,
  height,
  colSpan = 9,
}: Props) {
  if (emptyRows > 0) {
    return null;
  }

  return (
    <TableRow
      sx={{
        ...(height && {
          height: height * emptyRows,
        }),
      }}
    >
      <TableCell colSpan={colSpan}>No data available in table.</TableCell>
    </TableRow>
  );
}
