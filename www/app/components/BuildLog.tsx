/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
import Loading from './Loading';

interface OnClose {
	(): void;
}

interface Props {
	id: string;
	name: string;
	shown: boolean;
	onClose: OnClose;
}

interface State {
	closing: boolean;
}

const css = {
	loading: {
		float: "left",
	} as React.CSSProperties,
	buildLog: {
		top: '20px',
		width: 'calc(100% - 40px)',
		height: 'calc(100% - 40px)',
	} as React.CSSProperties,
	buildLogOutput: {
		fontSize: '10px',
		width: '100%',
		height: 'calc(100% - 130px)',
		overflow: 'scroll',
		padding: '2px 6px',
		backgroundColor: 'rgba(0, 0, 0, 0.2)',
	} as React.CSSProperties,
};

export default class Build extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			closing: false,
		};
	}

	closeDialog = (): void => {
		this.setState({
			closing: true,
		});
		this.props.onClose();
		setTimeout(() => {
			this.setState({
				closing: false,
			});
		}, 500);
	}

	render(): JSX.Element {
		if (!this.props.shown && !this.state.closing) {
			return null;
		}

		return <Blueprint.Dialog
			title={`Builds Logs - ${this.props.name}`}
			style={css.buildLog}
			isOpen={this.props.shown}
			onClose={this.closeDialog}
			canOutsideClickClose={false}
			>
			<div className="pt-dialog-body">
					<pre style={css.buildLogOutput}>
						{['test'] ? ['test'].join('\n') : ''}
					</pre>
			</div>
			<div className="pt-dialog-footer">
				<Loading size="small" style={css.loading}/>
				<div className="pt-dialog-footer-actions">
					<button type="button"
						className="pt-button"
						onClick={this.closeDialog}
					>Close</button>
				</div>
			</div>
		</Blueprint.Dialog>;
	}
}
