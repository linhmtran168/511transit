export type Operator = {
  name: string;
  id: string;
}

export type MainStore = {
  operators: Operator[];
  tripUpdates: TripUpdateEntry[];
  selectedOperator: string | null;
  isWaitingTripUpdates: boolean;
}

export enum MessageType {
  operators = 'operators',
  tripUpdates = 'tripUpdates'
}

export type TripUpdateParm = {
  operatorId: string;
}

export type OperatorsRequest = {
  requestType: MessageType;
}

export type TripUpdatesRequest = {
  requestType: MessageType;
  data: TripUpdateParm;
}

export type MessageRequest = OperatorsRequest | TripUpdatesRequest

export type OperatorsResponse = {
  responseType: MessageType;
  data: Operator[];
}

export type TripUpdatesResponse = {
  responseType: MessageType;
  data: any;
}

export type MessageResponse = OperatorsResponse | TripUpdatesResponse

export type TripUpdatesData = {
  operatorId: string;
  tripUpdates: TripUpdateEntry[];
}

export type TripUpdateEntry = {
  id: string;
  tripUpdate: TripUpdate;
}

export type TripUpdate = {
  trip: Trip;
  vehicle: Vehicle;
  stopTimeUpdate: StopTimeEntry[];
  timestamp: number;
}


export type StopTime = { time: number; delay?: number; uncertainty?: number }
export type StopTimeEntry = {
  arrival?: StopTime;
  departure?: StopTime;
  scheduleRelationship?: number;
  stopSequence: number;
  stopId: string;
}

export type Trip = {
  tripId: string;
  routeId: string;
  directionId: number;
  scheduleRelationship?: number;
}

export type Vehicle = {
  id: string;
  label: string;
  licensePlate: string;
}