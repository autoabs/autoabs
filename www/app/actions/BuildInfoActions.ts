/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as BuildInfoTypes from '../types/BuildInfoTypes';
import BuildInfoStore from '../stores/BuildInfoStore';
import * as MiscUtils from '../utils/MiscUtils';

let syncId: string;

export function sync(): Promise<void> {
	if (!BuildInfoStore.id) {
		return Promise.resolve();
	}

	let curSyncId = MiscUtils.uuid();
	syncId = curSyncId;

	let loader = new Loader().loading();

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/build/' + BuildInfoStore.id + '/log')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (curSyncId !== syncId) {
					resolve();
					return;
				}

				if (err) {
					Alert.error('Failed to get build log');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: BuildInfoTypes.UPDATE,
					data: {
						output: res.body,
					},
				});

				resolve();
			});
	});
}

export function open(id: string): Promise<void> {
	Dispatcher.dispatch({
		type: BuildInfoTypes.OPEN,
		data: {
			id: id,
		},
	});

	return sync();
}

export function close(): void {
	Dispatcher.dispatch({
		type: BuildInfoTypes.CLOSE,
	});
}
