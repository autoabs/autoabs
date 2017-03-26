/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import EventDispatcher from '../dispatcher/EventDispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as NodeTypes from '../types/NodeTypes';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/node')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.error('Failed to sync nodes');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: NodeTypes.SYNC,
					data: {
						nodes: res.body,
					},
				});

				resolve();
			});
	});
}

EventDispatcher.register((action: NodeTypes.NodeDispatch) => {
	switch (action.type) {
		case NodeTypes.CHANGE:
			sync();
			break;
	}
});
