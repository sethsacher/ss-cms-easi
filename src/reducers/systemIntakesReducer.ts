import { SystemIntakesState } from 'types/systemIntake';
import { prepareSystemIntakeForApp } from 'data/systemIntake';
import { fetchSystemIntakes } from 'routines/routines';
import { DateTime } from 'luxon';

const initialState: SystemIntakesState = {
  systemIntakes: []
};

function systemIntakesReducer(
  state = initialState,
  action: any
): SystemIntakesState {
  switch (action.type) {
    case fetchSystemIntakes.REQUEST:
      return {
        ...state,
        isLoading: true
      };
    case fetchSystemIntakes.SUCCESS:
      return {
        ...state,
        systemIntakes: action.payload.map((intake: any) =>
          prepareSystemIntakeForApp(intake)
        ),
        loadedTimestamp: DateTime.local()
      };
    case fetchSystemIntakes.FAILURE:
      return {
        ...state,
        error: action.payload
      };
    case fetchSystemIntakes.FULFILL:
      return {
        ...state,
        isLoading: false
      };
    default:
      return state;
  }
}

export default systemIntakesReducer;
