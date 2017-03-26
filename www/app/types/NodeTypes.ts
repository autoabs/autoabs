/// <reference path="../References.d.ts"/>
export const SYNC = 'node.sync';
export const CHANGE = 'node.change';

export interface Node {
	id: string;
	type: string;
	memory: number;
	load1: number;
	load5: number;
	load15: number;
}

export type Nodes = Node[];

export interface NodeDispatch {
	type: string;
	data?: {
		nodes?: Node[];
	};
}
