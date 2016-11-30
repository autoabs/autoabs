/// <reference path="../References.d.ts"/>
import * as React from 'react';
import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import LinearProgress from 'material-ui/LinearProgress';
import * as Constants from '../Constants';
import Styles from '../Styles';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	style?: React.CSSProperties;
	label?: string;
	primary?: boolean;
	disabled?: boolean;
	progressColor?: string;
	onConfirm?: () => void;
}

interface State {
	dialog: boolean;
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
	dialogOk: {
		color: Styles.colors.blue500,
	} as React.CSSProperties,
	dialogCancel: {
		color: Styles.colors.grey500,
	} as React.CSSProperties,
};

export default class ConfirmButton extends React.Component<Props, State> {
	constructor(props: Props, context: any) {
		super(props, context);
		this.state = {
			dialog: false,
			confirm: 0,
			confirming: null,
		};
	}

	openDialog = (): void => {
		this.setState(Object.assign({}, this.state, {
			dialog: true,
		}));
	}

	closeDialog = (): void => {
		this.setState(Object.assign({}, this.state, {
			dialog: false,
		}));
	}

	closeDialogConfirm = (): void => {
		this.setState(Object.assign({}, this.state, {
			dialog: false,
		}));
		if (this.props.onConfirm) {
			this.props.onConfirm();
		}
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
		let confirmElem: JSX.Element;

		if (Constants.mobile) {
			label = this.props.label;
			confirmElem = <Dialog
				title="Confirm"
				modal={true}
				open={this.state.dialog}
				actions={[
					<FlatButton
						label="Cancel"
						style={css.dialogCancel}
						onTouchTap={this.closeDialog}
					/>,
					<FlatButton
						label="Ok"
						style={css.dialogOk}
						onTouchTap={this.closeDialogConfirm}
					/>,
				]}
			/>;
		} else {
			if (this.state.confirming) {
				label = 'Hold';
				confirmElem = <LinearProgress
					style={css.actionProgress}
					color={this.props.progressColor}
					mode="determinate" max={10}
					value={this.state.confirm}
				/>;
			} else {
				label = this.props.label;
			}
		}

		return <div style={css.box}>
			<FlatButton
				style={this.props.style}
				label={label}
				primary={this.props.primary}
				disabled={this.props.disabled}
				onMouseDown={Constants.mobile ? undefined : this.confirm}
				onMouseUp={Constants.mobile ? undefined : this.clearConfirm}
				onMouseLeave={Constants.mobile ? undefined : this.clearConfirm}
				onTouchTap={Constants.mobile ? this.openDialog : undefined}
			/>
			{confirmElem}
		</div>;
	}
}
