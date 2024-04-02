import {
  TableEmptyRows,
  TableHeadCustom,
  TableNoData,
  TablePaginationCustom,
  TableSkeleton,
  useTable,
} from "@/components/table";
import { Card, Table, TableBody, TableContainer } from "@mui/material";
import React, { useCallback, useMemo, useState } from "react";
import Scrollbar from "@/components/scrollbar";
import VehicleTableRow from "./table-row";
import { DATA, TABLE_HEAD } from "./const";
import { useBoolean } from "@/hooks/use-boolean";
import _ from "lodash";

const defaultFilters = {
  name: "",
};
export default function VehicleTransectionTable() {
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const [filters, setFilters] = useState(defaultFilters);
  const dataInTable = useMemo(() => DATA, []);
  const handleDeleteRow = useCallback((id: string) => {
    //
  }, []);
  const openVehicleForm = useBoolean();
  const [editVehicle, setEditVehicle] = useState(undefined);
  const handleEdit = useCallback(
    (item: any) => {
      setEditVehicle(item);
      openVehicleForm.onTrue();
    },
    [openVehicleForm]
  );
  return (
    <Card>
      <TableContainer sx={{ position: "relative", overflow: "unset" }}>
        <Scrollbar>
          <Table size={"small"} sx={{ minWidth: 960 }}>
            <TableHeadCustom
              order={table.order}
              orderBy={table.orderBy}
              headLabel={TABLE_HEAD}
              rowCount={0}
              numSelected={table.selected.length}
              onSort={table.onSort}
            />
            {false ? (
              <TableSkeleton />
            ) : (
              <TableBody>
                {dataInTable
                  .slice(
                    table.page * table.rowsPerPage,
                    table.page * table.rowsPerPage + table.rowsPerPage
                  )
                  .map((row) => (
                    <VehicleTableRow
                      key={row.id}
                      row={row}
                      onDeleteRow={() => handleDeleteRow(row.id)}
                      onEdit={() => handleEdit(row)}
                    />
                  ))}
                <TableEmptyRows height={76} emptyRows={dataInTable.length} />
                <TableNoData
                  notFound={
                    !!Object.values(filters).filter((a) => a).length &&
                    dataInTable.length === 0
                  }
                />
              </TableBody>
            )}
          </Table>
        </Scrollbar>
      </TableContainer>

      <TablePaginationCustom
        disabled={false}
        count={dataInTable.length}
        page={table.page}
        rowsPerPage={table.rowsPerPage}
        onPageChange={table.onChangePage}
        onRowsPerPageChange={table.onChangeRowsPerPage}
        //
        dense={table.dense}
        onChangeDense={table.onChangeDense}
      />
    </Card>
  );
}
