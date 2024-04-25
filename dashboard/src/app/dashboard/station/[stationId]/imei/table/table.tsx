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
import Scrollbar from "@/components/scrollbar";
import ImeiTableRow from "./table-row";
import { TABLE_HEAD } from "./const";
import { useBoolean } from "@/hooks/use-boolean";
import _ from "lodash";
import { PropWithStationLocationId } from "../../types";
import { useGetStationImeiImsiActivitiesQuery } from "@/apollo-client";
import { transformData } from "./types";

const defaultFilters = {
  name: "",
};
export default function ImeiTransectionTable({
  stationLocationId,
}: PropWithStationLocationId) {
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const { data, loading } = useGetStationImeiImsiActivitiesQuery({
    variables: {
      limit: 1000,
      skip: 0,
      stationId: stationLocationId,
    },
  });
  const [filters, setFilters] = useState(defaultFilters);
  const dataInTable = useMemo(
    () =>
      data?.getStationImeiImsiActivities?.map((row) => transformData(row)) ??
      [],
    [data?.getStationImeiImsiActivities]
  );
  const handleDeleteRow = useCallback((id: string) => {
    //
  }, []);
  const openImeiForm = useBoolean();
  const [editImei, setEditImei] = useState(undefined);
  const handleEdit = useCallback(
    (item: any) => {
      setEditImei(item);
      openImeiForm.onTrue();
    },
    [openImeiForm]
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
                    <ImeiTableRow
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
