import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"


export class CreateScheduleRequest extends jspb.Message {
  getSchedulesList(): Array<Schedule>;
  setSchedulesList(value: Array<Schedule>): CreateScheduleRequest;
  clearSchedulesList(): CreateScheduleRequest;
  addSchedules(value?: Schedule, index?: number): Schedule;

  getForce(): boolean;
  setForce(value: boolean): CreateScheduleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateScheduleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateScheduleRequest): CreateScheduleRequest.AsObject;
  static serializeBinaryToWriter(message: CreateScheduleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateScheduleRequest;
  static deserializeBinaryFromReader(message: CreateScheduleRequest, reader: jspb.BinaryReader): CreateScheduleRequest;
}

export namespace CreateScheduleRequest {
  export type AsObject = {
    schedulesList: Array<Schedule.AsObject>,
    force: boolean,
  }
}

export class GetScheduleRequest extends jspb.Message {
  getScheduleId(): string;
  setScheduleId(value: string): GetScheduleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetScheduleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetScheduleRequest): GetScheduleRequest.AsObject;
  static serializeBinaryToWriter(message: GetScheduleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetScheduleRequest;
  static deserializeBinaryFromReader(message: GetScheduleRequest, reader: jspb.BinaryReader): GetScheduleRequest;
}

export namespace GetScheduleRequest {
  export type AsObject = {
    scheduleId: string,
  }
}

export class ListScheduleRequest extends jspb.Message {
  getScheduleIdsList(): Array<string>;
  setScheduleIdsList(value: Array<string>): ListScheduleRequest;
  clearScheduleIdsList(): ListScheduleRequest;
  addScheduleIds(value: string, index?: number): ListScheduleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListScheduleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListScheduleRequest): ListScheduleRequest.AsObject;
  static serializeBinaryToWriter(message: ListScheduleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListScheduleRequest;
  static deserializeBinaryFromReader(message: ListScheduleRequest, reader: jspb.BinaryReader): ListScheduleRequest;
}

export namespace ListScheduleRequest {
  export type AsObject = {
    scheduleIdsList: Array<string>,
  }
}

export class ListScheduleResponse extends jspb.Message {
  getSchedulesList(): Array<Schedule>;
  setSchedulesList(value: Array<Schedule>): ListScheduleResponse;
  clearSchedulesList(): ListScheduleResponse;
  addSchedules(value?: Schedule, index?: number): Schedule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListScheduleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListScheduleResponse): ListScheduleResponse.AsObject;
  static serializeBinaryToWriter(message: ListScheduleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListScheduleResponse;
  static deserializeBinaryFromReader(message: ListScheduleResponse, reader: jspb.BinaryReader): ListScheduleResponse;
}

export namespace ListScheduleResponse {
  export type AsObject = {
    schedulesList: Array<Schedule.AsObject>,
  }
}

export class UpdateScheduleRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateScheduleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateScheduleRequest): UpdateScheduleRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateScheduleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateScheduleRequest;
  static deserializeBinaryFromReader(message: UpdateScheduleRequest, reader: jspb.BinaryReader): UpdateScheduleRequest;
}

export namespace UpdateScheduleRequest {
  export type AsObject = {
  }
}

export class DeleteScheduleRequest extends jspb.Message {
  getScheduleId(): string;
  setScheduleId(value: string): DeleteScheduleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteScheduleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteScheduleRequest): DeleteScheduleRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteScheduleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteScheduleRequest;
  static deserializeBinaryFromReader(message: DeleteScheduleRequest, reader: jspb.BinaryReader): DeleteScheduleRequest;
}

export namespace DeleteScheduleRequest {
  export type AsObject = {
    scheduleId: string,
  }
}

export class Schedule extends jspb.Message {
  getId(): string;
  setId(value: string): Schedule;

  getDepartamentId(): string;
  setDepartamentId(value: string): Schedule;

  getDayOfWeek(): number;
  setDayOfWeek(value: number): Schedule;

  getHour(): number;
  setHour(value: number): Schedule;

  getMinute(): number;
  setMinute(value: number): Schedule;

  getVideoId(): string;
  setVideoId(value: string): Schedule;

  getCreateAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreateAt(value?: google_protobuf_timestamp_pb.Timestamp): Schedule;
  hasCreateAt(): boolean;
  clearCreateAt(): Schedule;

  getUpdateAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdateAt(value?: google_protobuf_timestamp_pb.Timestamp): Schedule;
  hasUpdateAt(): boolean;
  clearUpdateAt(): Schedule;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Schedule.AsObject;
  static toObject(includeInstance: boolean, msg: Schedule): Schedule.AsObject;
  static serializeBinaryToWriter(message: Schedule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Schedule;
  static deserializeBinaryFromReader(message: Schedule, reader: jspb.BinaryReader): Schedule;
}

export namespace Schedule {
  export type AsObject = {
    id: string,
    departamentId: string,
    dayOfWeek: number,
    hour: number,
    minute: number,
    videoId: string,
    createAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updateAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

