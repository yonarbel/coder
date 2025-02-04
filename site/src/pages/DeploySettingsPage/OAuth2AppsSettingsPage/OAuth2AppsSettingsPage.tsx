import { useQuery } from "react-query";
import { getApps } from "api/queries/oauth2";
import { useDashboard } from "components/Dashboard/DashboardProvider";
import { FC } from "react";
import { Helmet } from "react-helmet-async";
import { pageTitle } from "utils/page";
import OAuth2AppsSettingsPageView from "./OAuth2AppsSettingsPageView";

const OAuth2AppsSettingsPage: FC = () => {
  const { entitlements } = useDashboard();
  const appsQuery = useQuery(getApps());

  return (
    <>
      <Helmet>
        <title>{pageTitle("OAuth2 Applications")}</title>
      </Helmet>
      <OAuth2AppsSettingsPageView
        apps={appsQuery.data}
        isLoading={appsQuery.isLoading}
        isEntitled={
          entitlements.features.oauth2_provider.entitlement !== "not_entitled"
        }
      />
    </>
  );
};

export default OAuth2AppsSettingsPage;
