import ImeiView from "./view";

type Props = {
  params: { stationId: string };
};

export default function Page({ params }: Props) {
  return <ImeiView stationLocationId={params.stationId} />;
}
