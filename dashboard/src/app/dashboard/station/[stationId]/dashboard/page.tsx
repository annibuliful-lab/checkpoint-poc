import VehicleView from "./view";

type Props = {
  params: { stationId: string };
};

export default function Page({ params }: Props) {
  return <VehicleView stationLocationId={params.stationId} />;
}
