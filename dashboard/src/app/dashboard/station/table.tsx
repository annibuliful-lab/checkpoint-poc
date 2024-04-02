import {
  TableEmptyRows,
  TableHeadCustom,
  TableNoData,
  TablePaginationCustom,
  TableSkeleton,
  emptyRows,
  useTable,
} from "@/components/table";
import { Card, Table, TableBody, TableContainer } from "@mui/material";
import React, { useCallback, useMemo, useState } from "react";
import { IStationTableFilterValue, IStationTableFilters } from "./types";
import Scrollbar from "@/components/scrollbar";
import StationTableRow from "./table-row";
import { TABLE_HEAD } from "./const";
import {
  StationLocation,
  useDeleteStationLocationMutation,
  useGetStationLocationsQuery,
} from "@/apollo-client";
import { useBoolean } from "@/hooks/use-boolean";
import StationForm from "./station-form";
import _ from "lodash";

const defaultFilters: IStationTableFilters = {
  name: "",
};
export default function StationTable() {
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const [callDelete, callDeleteResponse] = useDeleteStationLocationMutation();
  const { loading, data } = useGetStationLocationsQuery({
    variables: {
      limit: 100,
      skip: 0,
    },
  });
  const [filters, setFilters] = useState(defaultFilters);
  const dataInTable = useMemo(
    () => (data?.getStationLocations ?? []) as StationLocation[],
    [data?.getStationLocations]
  );
  const handleFilters = useCallback(
    (name: string, value: IStationTableFilterValue) => {
      table.onResetPage();
      setFilters((prevState) => ({
        ...prevState,
        [name]: value,
      }));
    },
    [table]
  );
  const handleDeleteRow = useCallback(
    (id: string) => {
      callDelete({ variables: { id } });
    },
    [callDelete]
  );
  const openStationForm = useBoolean();
  const [editStation, setEditStation] = useState<StationLocation | undefined>(
    undefined
  );
  const handleEdit = useCallback(
    (item: StationLocation) => {
      setEditStation(item);
      openStationForm.onTrue();
    },
    [openStationForm]
  );
  return (
    <Card>
      <StationForm
        station={editStation}
        opened={openStationForm.value}
        onClose={() => {
          openStationForm.onFalse();
          setEditStation(undefined);
        }}
      />

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
            {loading ? (
              <TableSkeleton />
            ) : (
              <TableBody>
                {dataInTable
                  .slice(
                    table.page * table.rowsPerPage,
                    table.page * table.rowsPerPage + table.rowsPerPage
                  )
                  .map((row) => (
                    <StationTableRow
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
        disabled={loading}
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
