package aloha;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;

public class HttpConn {

    public class HttpAnswer {
        public int code;
        public String answer;
    }

    public HttpAnswer get(String GET_URL, String USER_AGENT) throws IOException {
        HttpAnswer ha = new HttpAnswer();
        URL obj = new URL(GET_URL);
        HttpURLConnection con = (HttpURLConnection) obj.openConnection();
        con.setRequestMethod("GET");
        con.setRequestProperty("User-Agent", USER_AGENT);
        ha.code = con.getResponseCode();

        StringBuffer response = new StringBuffer();
        if (ha.code == HttpURLConnection.HTTP_OK) { // success
            BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()));
            String inputLine;

            while ((inputLine = in.readLine()) != null) {
                response.append(inputLine);
            }
            in.close();

            ha.answer = response.toString();
        } else {
            ha.answer = "";
        }
        return ha;
    }

    public HttpAnswer post(String POST_URL, String USER_AGENT, String POST_PARAMS) throws IOException {
        HttpAnswer ha = new HttpAnswer();
        URL obj = new URL(POST_URL);
        HttpURLConnection con = (HttpURLConnection) obj.openConnection();
        con.setRequestMethod("POST");
        con.setRequestProperty("User-Agent", USER_AGENT);

        // For POST only - START
        con.setDoOutput(true);
        OutputStream os = con.getOutputStream();
        os.write(POST_PARAMS.getBytes());
        os.flush();
        os.close();
        // For POST only - END

        ha.code = con.getResponseCode();
        StringBuffer response = new StringBuffer();
        if (ha.code == HttpURLConnection.HTTP_OK) { // success
            BufferedReader in = new BufferedReader(new InputStreamReader(con.getInputStream()));
            String inputLine;

            while ((inputLine = in.readLine()) != null) {
                response.append(inputLine);
            }
            in.close();

            ha.answer = response.toString();
        } else {
            ha.answer = "";
        }
        return ha;
    }

}