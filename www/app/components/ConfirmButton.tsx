/// <reference path="../References.d.ts"/>
import * as React from 'react';
import FlatButton from 'material-ui/FlatButton';
import LinearProgress from 'material-ui/LinearProgress';
import Styles from '../Styles';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	style?: React.CSSProperties,
	label?: string,
	disabled?: boolean,
	progressColor?: string,
	onConfirm?: () => void,
}

interface State {
	confirm: number;
	confirming: string;
}

const css = {
	box: {
		display: 'inline-block',
		position: 'relative',
	} as React.CSSProperties,
	actionProgress: {
		position: 'absolute',
		bottom: 0,
		borderRadius: 0,
		borderBottomLeftRadius: '2px',
		borderBottomRightRadius: '2px',
		width: '100%',
		backgroundColor: 'rgba(0, 0, 0, 0)',
	} as React.CSSProperties,
};

export default class ConfirmButton extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			confirm: 0,
			confirming: null,
		};
	}

	confirm = (): void => {
		let confirmId = MiscUtils.uuid();

		this.setState(Object.assign({}, this.state, {
			confirming: confirmId,
		}));

		let i = 10;
		let id = setInterval(() => {
			if (i > 100) {
				clearInterval(id);
				setTimeout(() => {
					if (this.state.confirming === confirmId) {
						this.setState(Object.assign({}, this.state, {
							confirm: 0,
							confirming: null,
						}));
						if (this.props.onConfirm) {
							this.props.onConfirm();
						}
					}
				}, 365);
				return;
			} else if (!this.state.confirming) {
				clearInterval(id);
				this.setState(Object.assign({}, this.state, {
					confirm: 0,
					confirming: null,
				}));
				return;
			}

			if (i % 10 === 0) {
				this.setState(Object.assign({}, this.state, {
					confirm: i / 10,
				}));
			}

			i += 1;
		}, 8);
	}

	clearConfirm = (): void => {
		this.setState(Object.assign({}, this.state, {
			confirm: 0,
			confirming: null,
		}));
	}

	render(): JSX.Element {
		let label: string;
		let progress: JSX.Element;

		if (this.state.confirming) {
			label = 'Hold';
			progress = <LinearProgress
				style={css.actionProgress}
				color={this.props.progressColor}
				mode="determinate" max={10}
				value={this.state.confirm}
			/>;
		} else {
			label = this.props.label;
		}

		return <div style={css.box}>
			<FlatButton
				style={this.props.style}
				label={label}
				disabled={this.props.disabled}
				primary={true}
				onMouseDown={this.confirm}
				onMouseUp={this.clearConfirm}
			/>
			{progress}
		</div>;
	}
}
